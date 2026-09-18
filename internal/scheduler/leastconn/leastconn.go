package leastconnbalancer

import (
	"container/heap"

	"usman-faisal/tcp-loadbalancer/internal/queue"
	"usman-faisal/tcp-loadbalancer/internal/scheduler/common"
)

func New(backendList []string) *Scheduler {
	s := &Scheduler{}
	heap.Init(&s.heap)
	for _, addr := range backendList {
		b := &Backend{
			Addr:      addr,
			IsHealthy: true,
			Queue:     queue.New(common.DEFAULT_QUEUE_SIZE),
			Sem:       make(chan struct{}, common.DEFAULT_SEM_SIZE),
		}
		heap.Push(&s.heap, b)
	}
	for _, b := range s.heap {
		go s.Drain(b)
	}
	return s
}
