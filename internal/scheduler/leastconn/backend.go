package leastconnbalancer

type Backend struct {
	Addr        string
	ActiveConns int
	index       int
	IsHealthy   bool
}

func (b *Backend) GetAddr() string {
	return b.Addr
}

func (b *Backend) GetHealth() bool {
	return b.IsHealthy
}
