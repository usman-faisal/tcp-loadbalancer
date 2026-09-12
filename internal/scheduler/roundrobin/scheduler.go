package roundrobin

import (
	"fmt"
	"sync"
	"sync/atomic"
	"usman-faisal/tcp-loadbalancer/internal/types"
)

type SafeRoundRobin struct {
	R  RoundRobin
	mu sync.RWMutex
}

func New(backendList []string) *SafeRoundRobin {
	if len(backendList) == 0 {
		return nil
	}

	var backends []*Backend
	for _, v := range backendList {
		backends = append(backends, &Backend{
			Addr:      v,
			IsHealthy: true,
		})

	}

	return &SafeRoundRobin{
		R: RoundRobin{
			backends: backends,
			index:    0,
		},
	}
}

func (rb *SafeRoundRobin) Pick() types.IsBackend {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	n := len(rb.R.backends)
	if n == 0 {
		return nil
	}

	idx := atomic.LoadUint32(&rb.R.index)
	for i := 0; i < n; i++ {
		b := rb.R.backends[(idx+uint32(i))%uint32(n)]
		if b.IsHealthy {
			return b
		}
	}
	return nil
}
func (rb *SafeRoundRobin) Handle(b types.IsBackend) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.R.Next()
}

func (rb *SafeRoundRobin) Cleanup(b types.IsBackend) {
	// todo:
	fmt.Printf("cleanUp")

}

func (rb *SafeRoundRobin) Snapshot() {
	for i, backend := range rb.R.backends {
		fmt.Printf("[%d] addr=%s",
			i, backend.Addr)
	}
}

func (rb *SafeRoundRobin) SetHealth(b types.IsBackend, status bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	backend, ok := b.(*Backend)
	if !ok {
		return
	}
	backend.IsHealthy = status
}
