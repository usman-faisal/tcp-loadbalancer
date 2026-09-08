package leastconnbalancer

import (
	"fmt"
	"sync"
)

type LeastConnBalancer struct {
	mu   sync.RWMutex
	Heap BackendHeap
}

func New(backendList []string) *LeastConnBalancer {
	leastConnBalancer := LeastConnBalancer{}
	for _, b := range backendList {
		backend := &Backend{
			Addr:        b,
			ActiveConns: 0,
		}

		leastConnBalancer.Heap.Push(backend)
	}
	return &leastConnBalancer
}

func (lc *LeastConnBalancer) Pick() *Backend {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if len(lc.Heap) == 0 {
		return nil
	}

	return lc.Heap[0]
}
func (lc *LeastConnBalancer) Handle(b *Backend) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.Heap.update(b, 1)
}
func (lc *LeastConnBalancer) Cleanup(b *Backend) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.Heap.update(b, -1)
}
func (lc *LeastConnBalancer) Snapshot() {
	for i, backend := range lc.Heap {
		fmt.Printf("[%d] addr=%s conns=%d heapIdx=%d\n",
			i, backend.Addr, backend.ActiveConns, backend.index)
	}
}
