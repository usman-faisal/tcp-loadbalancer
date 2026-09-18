package leastconnbalancer

import "container/heap"

type BackendHeap []*Backend

func (h BackendHeap) Len() int { return len(h) }

func (h BackendHeap) Less(i, j int) bool {
	if h[i].IsHealthy != h[j].IsHealthy {
		return h[i].IsHealthy
	}
	return h[i].ActiveConns < h[j].ActiveConns
}

func (h BackendHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *BackendHeap) Push(x any) {
	n := len(*h)
	b := x.(*Backend)
	b.Index = n
	*h = append(*h, b)
}

func (h *BackendHeap) Pop() any {
	old := *h
	n := len(old)
	b := old[n-1]
	old[n-1] = nil
	b.Index = -1
	*h = old[0 : n-1]
	return b
}

func (h *BackendHeap) update(b *Backend, delta int) {
	b.ActiveConns += delta
	heap.Fix(h, b.Index)
}
