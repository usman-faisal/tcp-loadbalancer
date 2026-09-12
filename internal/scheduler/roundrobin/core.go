package roundrobin

import (
	"sync/atomic"
)

type RoundRobin struct {
	backends []*Backend
	index    uint32
}

func (r *RoundRobin) Curr() *Backend {
	n := atomic.LoadUint32(&r.index)
	return r.backends[n%uint32(len(r.backends))]
}

func (r *RoundRobin) Next() *Backend {
	n := atomic.AddUint32(&r.index, 1)
	return r.backends[n%uint32(len(r.backends))]
}
