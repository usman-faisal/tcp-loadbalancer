package leastconnbalancer

import (
	"container/heap"
	"fmt"
	"sync"
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
		}

		heap.Push(&leastConnBalancer.Heap, backend)
	}
	return &leastConnBalancer
}

func (lc *LeastConnBalancer) Pick() types.IsBackend {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if len(lc.Heap) == 0 {
		return nil
	}

	best := lc.Heap[0]
	if !best.IsHealthy {
		return nil
	}
	return best
}

func (lc *LeastConnBalancer) Handle(b types.IsBackend) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	backend := b.(*Backend)
	lc.Heap.update(backend, 1)
}

func (lc *LeastConnBalancer) Cleanup(b types.IsBackend) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	backend := b.(*Backend)
	lc.Heap.update(backend, -1)
}

func (lc *LeastConnBalancer) Snapshot() {
	for i, backend := range lc.Heap {
		fmt.Printf("[%d] addr=%s conns=%d heapIdx=%d\n",
			i, backend.Addr, backend.ActiveConns, backend.index)
	}
}

func (lc *LeastConnBalancer) SetHealth(b types.IsBackend, status bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	backend := b.(*Backend)
	backend.IsHealthy = status

	if backend.IsHealthy != status {
		backend.IsHealthy = status
		heap.Fix(&lc.Heap, backend.index)
	}
}
