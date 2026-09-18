package roundrobin

import "usman-faisal/tcp-loadbalancer/internal/queue"

func New(backendList []string) *Scheduler {
	if len(backendList) == 0 {
		return nil
	}
	s := &Scheduler{}
	for _, addr := range backendList {
		s.backends = append(s.backends, &Backend{
			Addr:      addr,
			IsHealthy: true,
			Queue:     queue.New(5),
			Sem:       make(chan struct{}, 50),
		})
	}
	for _, b := range s.backends {
		go s.Drain(b)
	}
	return s
}
