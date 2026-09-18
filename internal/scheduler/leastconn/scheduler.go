package leastconnbalancer

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"usman-faisal/tcp-loadbalancer/internal/transport"
	"usman-faisal/tcp-loadbalancer/internal/types"
)

type Scheduler struct {
	mu   sync.RWMutex
	heap BackendHeap
}

func (s *Scheduler) pickAndReserve() *Backend {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.heap) == 0 || !s.heap[0].IsHealthy {
		return nil
	}
	best := s.heap[0]
	s.heap.update(best, 1)
	return best
}

func (s *Scheduler) release(b *Backend) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.heap.update(b, -1)
}

func (s *Scheduler) Submit(conn net.Conn) (types.IsBackend, error) {
	b := s.pickAndReserve()
	if b == nil {
		return nil, errors.New("no healthy backend")
	}

	select {
	case b.Sem <- struct{}{}:
		go s.handle(conn, b)
	default:
		if err := b.Queue.Enqueue(conn, 5*time.Second); err != nil {
			s.release(b)
			conn.Close()
			return nil, err
		}
	}

	return b, nil
}

func (s *Scheduler) Cleanup(b types.IsBackend) {
	s.release(b.(*Backend))
}

func (s *Scheduler) SetHealth(b types.IsBackend, status bool) {
	backend := b.(*Backend)
	s.mu.Lock()
	defer s.mu.Unlock()
	old := backend.IsHealthy
	backend.IsHealthy = status
	if old != status {
		heap.Fix(&s.heap, backend.Index)
	}
}

func (s *Scheduler) Snapshot() {
	s.mu.RLock()
	snapshot := make([]*Backend, len(s.heap))
	copy(snapshot, s.heap)
	s.mu.RUnlock()

	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].ActiveConns < snapshot[j].ActiveConns
	})

	var sb strings.Builder
	sb.WriteString("--- backend snapshot (leastconn) ---\n")
	for i, b := range snapshot {
		status := "UP"
		if !b.IsHealthy {
			status = "DOWN"
		}
		fmt.Fprintf(&sb, "[%d] addr=%s conns=%d queued=%d %s\n", i, b.Addr, b.ActiveConns, b.Queue.Len(), status)
	}
	log.Print(sb.String())
}

func (s *Scheduler) handle(conn net.Conn, b *Backend) {
	dialed, err := net.Dial("tcp", b.GetAddr())
	if err != nil {
		s.SetHealth(b, false)
		log.Println(err)
		s.release(b)
		<-b.Sem
		go func(c net.Conn) {
			if _, err := s.Submit(c); err != nil {
				c.Close()
			}
		}(conn)
		return
	}

	transport.Proxy(dialed, conn, func() {
		s.Snapshot()
		s.release(b)
		<-b.Sem
	})
}

func (s *Scheduler) Drain(b *Backend) {
	for conn := range b.Queue.Waiting {
		if !b.GetHealth() {
			s.release(b)

			// reroute to another backend
			go func(c net.Conn) {
				if _, err := s.Submit(c); err != nil {
					c.Close()
				}
			}(conn)
			continue
		}
		b.Sem <- struct{}{}
		go s.handle(conn, b)
	}
}
