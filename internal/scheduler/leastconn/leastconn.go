package leastconnbalancer

import (
	"container/heap"

	"usman-faisal/tcp-loadbalancer/internal/queue"
)

func New(backendList []string) *Scheduler {
	s := &Scheduler{}
	heap.Init(&s.heap)
	for _, addr := range backendList {
		b := &Backend{
			Addr:      addr,
			IsHealthy: true,
			Queue:     queue.New(5),
			Sem:       make(chan struct{}, 50),
		}
		heap.Push(&s.heap, b)
	}
	for _, b := range s.heap {
		go s.Drain(b)
	}
	return s
}
