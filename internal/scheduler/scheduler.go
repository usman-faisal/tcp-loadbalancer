package scheduler

import (
	"net"
	leastconnbalancer "usman-faisal/tcp-loadbalancer/internal/scheduler/leastconn"
	"usman-faisal/tcp-loadbalancer/internal/scheduler/roundrobin"
	"usman-faisal/tcp-loadbalancer/internal/types"
)

type Scheduler interface {
	// Pick() types.IsBackend
	// Handle(b types.IsBackend)
	Submit(conn net.Conn) (types.IsBackend, error)
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
