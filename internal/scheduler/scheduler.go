package scheduler

type Algorithm string

const (
	RoundRobin Algorithm = "round-robin"
	LeastConn  Algorithm = "least-conn"
)

type IsBackend interface {
	GetAddr() string
}

type Scheduler[T IsBackend] interface {
	Pick() T
	Handle(b T)
	Cleanup(b T)
	Snapshot()
}

