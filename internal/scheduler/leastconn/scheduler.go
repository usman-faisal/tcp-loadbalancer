package leastconnbalancer

import (
	"container/heap"
	"errors"
	"log"
	"net"
	"sort"
	"sync"
	"time"
	"usman-faisal/tcp-loadbalancer/internal/queue"
	"usman-faisal/tcp-loadbalancer/internal/transport"
	"usman-faisal/tcp-loadbalancer/internal/types"
)

type LeastConnBalancer struct {
	mu   sync.RWMutex
	Heap BackendHeap
}

func New(backendList []string) *LeastConnBalancer {
	leastConnBalancer := LeastConnBalancer{}

	heap.Init(&leastConnBalancer.Heap)
	for _, b := range backendList {
		backend := &Backend{
			Addr:        b,
			ActiveConns: 0,
			IsHealthy:   true,
			queue:       queue.New(5),
			sem:         make(chan struct{}, 50),
		}
		go leastConnBalancer.drain(backend)
		heap.Push(&leastConnBalancer.Heap, backend)
	}
	return &leastConnBalancer
}

func (lc *LeastConnBalancer) Submit(conn net.Conn) (types.IsBackend, error) {
	best := lc.PickAndReserve()
	if best == nil {
		return nil, errors.New("no healthy backend")
	}
	backend := best.(*Backend)

	select {
	case backend.sem <- struct{}{}:
		go lc.handle(conn, backend)
	default:
		if err := backend.queue.Enqueue(conn, time.Second*5); err != nil {
			lc.Cleanup(backend)
			conn.Close()
			return nil, err
		}
	}

	return backend, nil
}
func (lc *LeastConnBalancer) drain(b types.IsBackend) {
	backend := b.(*Backend)
	for conn := range backend.queue.Waiting {
		if !backend.GetHealth() {
			conn.Close()
			lc.Cleanup(backend)
			continue
		}

		backend.sem <- struct{}{}
		go lc.handle(conn, backend)
	}
}

func (lc *LeastConnBalancer) handle(conn net.Conn, b types.IsBackend) {
	backend := b.(*Backend)

	dialedBackend, err := net.Dial("tcp", backend.GetAddr())

	if err != nil {
		lc.SetHealth(backend, false)
		log.Println(err)
		lc.Cleanup(backend)
		conn.Close()
		<-backend.sem
	}

	transport.Proxy(dialedBackend, conn, func() {
		lc.Snapshot()
		lc.Cleanup(backend)
		<-backend.sem
	})
}

func (lc *LeastConnBalancer) PickAndReserve() types.IsBackend {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if len(lc.Heap) == 0 || !lc.Heap[0].IsHealthy {
		return nil
	}
	best := lc.Heap[0]
	lc.Heap.update(best, 1)
	return best
}

func (lc *LeastConnBalancer) Cleanup(b types.IsBackend) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	backend := b.(*Backend)
	lc.Heap.update(backend, -1)
}

func (lc *LeastConnBalancer) Snapshot() {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	snapshot := make([]*Backend, len(lc.Heap))
	copy(snapshot, lc.Heap)

	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].ActiveConns < snapshot[j].ActiveConns
	})

	log.Println("--- backend snapshot ---")
	for i, b := range snapshot {
		status := "UP"
		if !b.IsHealthy {
			status = "DOWN"
		}
		log.Printf("[%d] addr=%s conns=%d queued=%d healthy=%v",
			i, b.Addr, b.ActiveConns, b.queue.Len(), status)
	}
}

func (lc *LeastConnBalancer) SetHealth(b types.IsBackend, status bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	backend := b.(*Backend)
	oldStatus := backend.IsHealthy
	backend.IsHealthy = status
	if oldStatus != status {
		heap.Fix(&lc.Heap, backend.index)
	}
}
