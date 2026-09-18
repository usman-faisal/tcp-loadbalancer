package roundrobin

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"usman-faisal/tcp-loadbalancer/internal/queue"
	"usman-faisal/tcp-loadbalancer/internal/transport"
	"usman-faisal/tcp-loadbalancer/internal/types"
)

var ErrNoBackends = errors.New("roundrobin: no backends provided")

type SafeRoundRobin struct {
	mu sync.Mutex
	r  RoundRobin
}

func New(backendList []string) *SafeRoundRobin {
	if len(backendList) == 0 {
		return nil
	}

	safeRoundRobin := &SafeRoundRobin{
		r: RoundRobin{
			backends: make([]*Backend, 0, len(backendList)),
			index:    0,
		},
	}

	for _, addr := range backendList {
		safeRoundRobin.r.backends = append(safeRoundRobin.r.backends, &Backend{
			Addr:      addr,
			IsHealthy: true,
			queue:     queue.New(2),
		})
	}

	for _, backend := range safeRoundRobin.r.backends {
		go safeRoundRobin.ProcessRequestsPerBackend(backend)
	}

	return safeRoundRobin
}

func (rb *SafeRoundRobin) Pick() types.IsBackend {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	b := rb.r.pick()
	if b == nil {
		return nil
	}
	return b
}

func (rb *SafeRoundRobin) Submit(conn net.Conn) (types.IsBackend, error) {
	best := rb.Pick()
	if best == nil {
		return nil, errors.New("no healthy backend")
	}
	backend := best.(*Backend)

	if !backend.IsHealthy {
		return nil, errors.New("No healthy backend")
	}

	if err := backend.queue.Enqueue(conn, time.Second*5); err != nil {
		return nil, err
	}

	return backend, nil
}

func (rb *SafeRoundRobin) ProcessRequestsPerBackend(b types.IsBackend) {
	backend := b.(*Backend)
	for conn := range backend.queue.Waiting {
		time.Sleep(2 * time.Second)
		dialedBackend, err := net.Dial("tcp", backend.GetAddr())
		if err != nil {
			rb.SetHealth(backend, false)
			log.Println(err)
			rb.Cleanup(backend)
			conn.Close()
			continue
		}
		go transport.Proxy(dialedBackend, conn, func() {
			rb.Snapshot()
			rb.Cleanup(backend)
		})
	}
}

func (rb *SafeRoundRobin) Cleanup(b types.IsBackend) {}

func (rb *SafeRoundRobin) SetHealth(b types.IsBackend, status bool) {
	backend, ok := b.(*Backend)
	if !ok || backend == nil {
		return
	}

	rb.mu.Lock()
	defer rb.mu.Unlock()
	backend.IsHealthy = status
}

func (rb *SafeRoundRobin) Snapshot() {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	log.Println("--- backend snapshot ---")
	for _, backend := range rb.r.backends {
		status := "UP"
		if !backend.IsHealthy {
			status = "DOWN"
		}
		log.Printf("  %-20s queued=%-3d %s\n", backend.Addr, backend.queue.Len(), status)
	}
}
