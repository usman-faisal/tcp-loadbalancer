package common

import (
	"errors"
	"log"
	"net"
	"time"
	"usman-faisal/tcp-loadbalancer/internal/transport"
	"usman-faisal/tcp-loadbalancer/internal/types"
)

type BaseScheduler struct {
	picker Picker
}

func NewBase(picker Picker) *BaseScheduler {
	return &BaseScheduler{picker: picker}
}

func (s *BaseScheduler) Submit(conn net.Conn) (types.IsBackend, error) {
	backend := s.picker.PickAndReserve()
	if backend == nil {
		return nil, errors.New("no healthy backend")
	}

	select {
	case backend.Sem <- struct{}{}:
		go s.handle(conn, backend)
	default:
		if err := backend.Queue.Enqueue(conn, 5*time.Second); err != nil {
			s.Cleanup(backend)
			conn.Close()
			return nil, err
		}
	}

	return backend, nil
}

func (s *BaseScheduler) Drain(backend *Backend) {
	for conn := range backend.Queue.Waiting {
		if !backend.GetHealth() {
			conn.Close()
			s.Cleanup(backend)
			continue
		}
		backend.Sem <- struct{}{}
		go s.handle(conn, backend)
	}
}

func (s *BaseScheduler) handle(conn net.Conn, backend *Backend) {
	dialedBackend, err := net.Dial("tcp", backend.GetAddr())
	if err != nil {
		s.SetHealth(backend, false)
		log.Println(err)
		s.Cleanup(backend)
		conn.Close()
		<-backend.Sem
		return
	}

	transport.Proxy(dialedBackend, conn, func() {
		s.Snapshot()
		s.Cleanup(backend)
		<-backend.Sem
	})
}

func (s *BaseScheduler) Cleanup(b types.IsBackend) {
	backend := b.(*Backend)
	s.picker.Release(backend)
}

func (s *BaseScheduler) SetHealth(b types.IsBackend, status bool) {
	backend := b.(*Backend)
	s.picker.SetHealth(backend, status)
}

func (s *BaseScheduler) Snapshot() {
	s.picker.Snapshot()
}
