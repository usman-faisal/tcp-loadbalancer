package leastconnbalancer

import (
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
	for _, b := range backendList {
		backend := &Backend{
			Addr:        b,
			ActiveConns: 0,
		}

		leastConnBalancer.Heap.Push(backend)
	}
	return &leastConnBalancer
}

func (lc *LeastConnBalancer) Pick() types.IsBackend{
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if len(lc.Heap) == 0 {
		return nil
	}

	return lc.Heap[0]
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
