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
