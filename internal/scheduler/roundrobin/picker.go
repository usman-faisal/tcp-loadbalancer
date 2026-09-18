package roundrobin

import (
	"log"
	"sync"
	"usman-faisal/tcp-loadbalancer/internal/scheduler/common"
)

type RoundRobinPicker struct {
	mu       sync.Mutex
	backends []*common.Backend
	index    int
}

func (p *RoundRobinPicker) PickAndReserve() *common.Backend {
	p.mu.Lock()
	defer p.mu.Unlock()

	n := len(p.backends)
	for i := 0; i < n; i++ {
		b := p.backends[p.index]
		p.index = (p.index + 1) % n
		if b.IsHealthy {
			b.ActiveConns++
			return b
		}
	}
	return nil
}

func (p *RoundRobinPicker) Release(b *common.Backend) {
	p.mu.Lock()
	defer p.mu.Unlock()
	b.ActiveConns--
}

func (p *RoundRobinPicker) SetHealth(b *common.Backend, status bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	b.IsHealthy = status
}

func (p *RoundRobinPicker) Snapshot() {
	p.mu.Lock()
	defer p.mu.Unlock()
	log.Println("--- backend snapshot (roundrobin) ---")
	for _, b := range p.backends {
		status := "UP"
		if !b.IsHealthy {
			status = "DOWN"
		}
		log.Printf("  %-20s conns=%-3d queued=%-3d %s\n", b.Addr, b.ActiveConns, b.Queue.Len(), status)
	}
}
