package roundrobin

import (
	"fmt"
	"sync"
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
		backends=append(backends, &Backend{
			Addr: v,
		})

	}

	return &SafeRoundRobin{
		R: RoundRobin{
			backends: backends,
			index: 0,
		},
	}
}

func (rb *SafeRoundRobin) Pick() types.IsBackend{
	rb.mu.Lock()
	defer rb.mu.Unlock()

	return rb.R.Curr()
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
