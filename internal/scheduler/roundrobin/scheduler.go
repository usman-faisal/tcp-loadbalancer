package roundrobin

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"usman-faisal/tcp-loadbalancer/internal/transport"
	"usman-faisal/tcp-loadbalancer/internal/types"
)

type Scheduler struct {
	mu       sync.Mutex
	backends []*Backend
	index    int
}

func (s *Scheduler) next() *Backend {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := len(s.backends)
	for i := 0; i < n; i++ {
		b := s.backends[s.index]
		s.index = (s.index + 1) % n
		if b.IsHealthy {
			return b
		}
	}
	return nil
}

func (s *Scheduler) Submit(conn net.Conn) (types.IsBackend, error) {
	b := s.next()
	if b == nil {
		return nil, errors.New("no healthy backend")
	}

	select {
	case b.Sem <- struct{}{}:
		go s.handle(conn, b)
	default:
		if err := b.Queue.Enqueue(conn, 5*time.Second); err != nil {
			conn.Close()
			return nil, err
		}
	}

	return b, nil
}

func (s *Scheduler) Cleanup(_ types.IsBackend) {}

func (s *Scheduler) SetHealth(b types.IsBackend, status bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b.(*Backend).IsHealthy = status
}

func (s *Scheduler) Snapshot() {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Println("--- backend snapshot (roundrobin) ---")
	for _, b := range s.backends {
		status := "UP"
		if !b.IsHealthy {
			status = "DOWN"
		}
		log.Printf("  %-20s queued=%-3d %s\n", b.Addr, b.Queue.Len(), status)
	}
}

func (s *Scheduler) handle(conn net.Conn, b *Backend) {
	dialed, err := net.Dial("tcp", b.GetAddr())
	if err != nil {
		s.SetHealth(b, false)
		log.Println(err)
		conn.Close()
		<-b.Sem
		return
	}

	transport.Proxy(dialed, conn, func() {
		s.Snapshot()
		<-b.Sem
	})
}

func (s *Scheduler) Drain(b *Backend) {
	for conn := range b.Queue.Waiting {
		if !b.GetHealth() {
			conn.Close()
			continue
		}
		b.Sem <- struct{}{}
		go s.handle(conn, b)
	}
}
