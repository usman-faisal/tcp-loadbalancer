package queue

import (
	"errors"
	"net"
	"time"
)

type ConnQueue struct {
	Waiting chan net.Conn
}

func New(maxSize int) *ConnQueue {
	return &ConnQueue{
		Waiting: make(chan net.Conn, maxSize),
	}
}

func (q *ConnQueue) Enqueue(conn net.Conn, timeout time.Duration) error {
	select {
	case q.Waiting <- conn:
		return nil
	case <-time.After(timeout):
		return errors.New("queue full, connection dropped")
	}
}

func (q *ConnQueue) Dequeue() (net.Conn, bool) {
	conn, ok := <-q.Waiting
	return conn, ok
}

func (q *ConnQueue) Len() int {
	return len(q.Waiting)
}
