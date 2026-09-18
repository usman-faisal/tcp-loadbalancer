package types

import "time"

type Health struct {
	IsHealthy   bool
	Lastchecked time.Time
}

type IsBackend interface {
	GetAddr() string
	GetHealth() bool
}

type Algorithm string

const (
	RoundRobin Algorithm = "round-robin"
	LeastConn  Algorithm = "least-conn"
)
type Stats struct {
    CPUPercent      float64 `json:"cpu_percent"`
    MemUsedBytes    uint64  `json:"mem_used_bytes"`
    MemTotalBytes   uint64  `json:"mem_total_bytes"`
    ActiveConns     int     `json:"active_connections"`
    // QueueDepth      int     `json:"queue_depth,omitempty"`
    Timestamp       int64   `json:"timestamp"`
}
