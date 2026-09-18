package roundrobin

import "usman-faisal/tcp-loadbalancer/internal/queue"

type Backend struct {
	Addr      string
	IsHealthy bool
	queue     *queue.ConnQueue
}

func (b *Backend) GetAddr() string { return b.Addr }

func (b *Backend) GetHealth() bool { return b.IsHealthy }
