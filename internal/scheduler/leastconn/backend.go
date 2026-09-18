package leastconnbalancer

import "usman-faisal/tcp-loadbalancer/internal/queue"

type Backend struct {
	Addr        string
	ActiveConns int
	index       int
	IsHealthy   bool
	queue       *queue.ConnQueue
	sem         chan struct{}
}

func (b *Backend) GetAddr() string {
	return b.Addr
}

func (b *Backend) GetHealth() bool {
	return b.IsHealthy
}
