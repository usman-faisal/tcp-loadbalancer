package roundrobin

import (
	"usman-faisal/tcp-loadbalancer/internal/queue"
	"usman-faisal/tcp-loadbalancer/internal/scheduler/common"
)

func New(backendList []string) *common.BaseScheduler {
	if len(backendList) == 0 {
		return nil
	}
	picker := &RoundRobinPicker{}
	for _, addr := range backendList {
		picker.backends = append(picker.backends, &common.Backend{
			Addr: addr, IsHealthy: true, Queue: queue.New(5), Sem: make(chan struct{}, 50),
		})
	}
	base := common.NewBase(picker)
	for _, b := range picker.backends {
		go base.Drain(b)
	}
	return base
}
