package leastconnbalancer

import (
	"container/heap"
	"usman-faisal/tcp-loadbalancer/internal/queue"
	"usman-faisal/tcp-loadbalancer/internal/scheduler/common"
)

func New(backendList []string) *common.BaseScheduler {
	picker := &LeastConnPicker{}
	heap.Init(&picker.Heap)
	for _, addr := range backendList {
		b := &common.Backend{Addr: addr, IsHealthy: true, Queue: queue.New(5), Sem: make(chan struct{}, 50)}
		heap.Push(&picker.Heap, b)
	}
	base := common.NewBase(picker)
	for _, b := range picker.Heap {
		go base.Drain(b)
	}
	return base
}
