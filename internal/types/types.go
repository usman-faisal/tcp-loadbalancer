package types

// IsBackend defines the contract for all backend structs
type IsBackend interface {
	GetAddr() string
}

type Algorithm string

const (
	RoundRobin Algorithm = "round-robin"
	LeastConn  Algorithm = "least-conn"
)
