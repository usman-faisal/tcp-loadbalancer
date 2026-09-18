package roundrobin

type RoundRobin struct {
	backends []*Backend
	index    int
}

func (r *RoundRobin) pick() *Backend {
	n := len(r.backends)
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		b := r.backends[r.index]
		r.index = (r.index + 1) % n
		if b.IsHealthy {
			return b
		}
	}
	return nil
}
