package leastconnbalancer

import (
	"container/heap"
)

type Backend struct {
	Addr string
	ActiveConns int
	index       int
}

type BackendHeap []*Backend

func (b *Backend) GetAddr() string {
	return b.Addr
}

func (pq BackendHeap) Len() int { return len(pq) }

func (pq BackendHeap) Less(i, j int) bool {
	return pq[i].ActiveConns < pq[j].ActiveConns
}

func (pq BackendHeap) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *BackendHeap) Push(x any) {
	n := len(*pq)
	backend := x.(*Backend)
	backend.index = n
	*pq = append(*pq, backend)
}

func (pq *BackendHeap) Pop() any {
	old := *pq
	n := len(old)
	backend := old[n-1]
	old[n-1] = nil
	backend.index = -1
	*pq = old[0 : n-1]
	return backend
}

func (pq *BackendHeap) update(backend *Backend, delta int) {
	backend.ActiveConns += delta
	heap.Fix(pq, backend.index)
}
