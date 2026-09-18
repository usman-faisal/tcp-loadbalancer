package roundrobin

import (
	"usman-faisal/tcp-loadbalancer/internal/queue"
	"usman-faisal/tcp-loadbalancer/internal/scheduler/common"
)

func New(backendList []string) *Scheduler {
	if len(backendList) == 0 {
		return nil
	}
	s := &Scheduler{}
	for _, addr := range backendList {
		s.backends = append(s.backends, &Backend{
			Addr:      addr,
			IsHealthy: true,
			Queue:     queue.New(common.DEFAULT_QUEUE_SIZE),
			Sem:       make(chan struct{}, common.DEFAULT_SEM_SIZE),
		})
	}
	for _, b := range s.backends {
		go s.Drain(b)
	}
	return s
}
