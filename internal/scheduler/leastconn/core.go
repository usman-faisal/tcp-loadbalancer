package leastconnbalancer

import (
	"container/heap"
	"usman-faisal/tcp-loadbalancer/internal/scheduler/common"
)

type BackendHeap []*common.Backend

func (pq BackendHeap) Len() int { return len(pq) }

func (pq BackendHeap) Less(i, j int) bool {
	if pq[i].IsHealthy != pq[j].IsHealthy {
		return pq[i].IsHealthy
	}
	return pq[i].ActiveConns < pq[j].ActiveConns
}

func (pq BackendHeap) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *BackendHeap) Push(x any) {
	n := len(*pq)
	backend := x.(*common.Backend)
	backend.Index = n
	*pq = append(*pq, backend)
}

func (pq *BackendHeap) Pop() any {
	old := *pq
	n := len(old)
	backend := old[n-1]
	old[n-1] = nil
	backend.Index = -1
	*pq = old[0 : n-1]
	return backend
}

func (pq *BackendHeap) update(backend *common.Backend, delta int) {
	backend.ActiveConns += delta
	heap.Fix(pq, backend.Index)
}
