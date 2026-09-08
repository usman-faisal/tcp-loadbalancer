package roundrobin

import (
	"fmt"
	"sync"
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

func (rb *SafeRoundRobin) Pick() *Backend {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	return rb.R.Curr()
}
func (rb *SafeRoundRobin) Handle(b *Backend) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.R.Next()
}

func (rb *SafeRoundRobin) Cleanup(b *Backend) {
	// todo:
	fmt.Printf("cleanUp")

}

func (rb *SafeRoundRobin) Snapshot() {
	for i, backend := range rb.R.backends {
		fmt.Printf("[%d] addr=%s",
			i, backend.Addr)
	}
}
