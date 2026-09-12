package scheduler

import (
	leastconnbalancer "usman-faisal/tcp-loadbalancer/internal/scheduler/leastconn"
	"usman-faisal/tcp-loadbalancer/internal/scheduler/roundrobin"
	"usman-faisal/tcp-loadbalancer/internal/types"
)

type Scheduler interface {
	Pick() types.IsBackend
	Handle(b types.IsBackend)
	Cleanup(b types.IsBackend)
	SetHealth(b types.IsBackend, health bool)
	Snapshot()
}

func Init(algorithm types.Algorithm, backendList []string) Scheduler {
	switch algorithm {
	case types.RoundRobin:
		return roundrobin.New(backendList)
	case types.LeastConn:
		return leastconnbalancer.New(backendList)
	default:
		return nil
	}
}
