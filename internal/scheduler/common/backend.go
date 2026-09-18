package common

import (
	"usman-faisal/tcp-loadbalancer/internal/queue"
)

type Backend struct {
	Addr        string
	ActiveConns int
	Index       int
	IsHealthy   bool
	Queue       *queue.ConnQueue
	Sem         chan struct{}
}

func (b *Backend) GetAddr() string { return b.Addr }
func (b *Backend) GetHealth() bool { return b.IsHealthy }
