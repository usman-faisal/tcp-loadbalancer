package leastconnbalancer

import (
	"container/heap"
	"log"
	"sort"
	"sync"
	"usman-faisal/tcp-loadbalancer/internal/scheduler/common"
)

type LeastConnPicker struct {
	mu   sync.RWMutex
	Heap BackendHeap
}

func (p *LeastConnPicker) PickAndReserve() *common.Backend {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.Heap) == 0 || !p.Heap[0].IsHealthy {
		return nil
	}
	best := p.Heap[0]
	p.Heap.update(best, 1)
	return best
}

func (p *LeastConnPicker) Release(b *common.Backend) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Heap.update(b, -1)
}

func (p *LeastConnPicker) SetHealth(b *common.Backend, status bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	old := b.IsHealthy
	b.IsHealthy = status
	if old != status {
		heap.Fix(&p.Heap, b.Index)
	}
}

func (p *LeastConnPicker) Snapshot() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	snapshot := make([]*common.Backend, len(p.Heap))
	copy(snapshot, p.Heap)
	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].ActiveConns < snapshot[j].ActiveConns
	})
	log.Println("--- backend snapshot (leastconn) ---")
	for i, b := range snapshot {
		status := "UP"
		if !b.IsHealthy {
			status = "DOWN"
		}
		log.Printf("[%d] addr=%s conns=%d queued=%d healthy=%v", i, b.Addr, b.ActiveConns, b.Queue.Len(), status)
	}
}
