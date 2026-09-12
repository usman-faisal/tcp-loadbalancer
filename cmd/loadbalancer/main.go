package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"usman-faisal/tcp-loadbalancer/internal/config"
	"usman-faisal/tcp-loadbalancer/internal/scheduler"
)

func proxy(backend net.Conn, conn net.Conn, cleanup func()) {
	defer backend.Close()
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(backend, conn)
		if tc, ok := backend.(*net.TCPConn); ok {
			tc.CloseWrite()
		}
	}()
	go func() {
		defer wg.Done()
		io.Copy(conn, backend)
		if tc, ok := conn.(*net.TCPConn); ok {
			tc.CloseWrite()
		}
	}()

	wg.Wait()

	cleanup()
}

func main() {
	cfg, err := config.Load("config.yaml")

	if err != nil {
		log.Println(err)
		return
	}

	backendList, port := cfg.BackendList, cfg.Port

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("listening on %s", ln.Addr().String())

	s := scheduler.Init(cfg.Algorithm, backendList)

	for {
		// listen for requests
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("accepting connection %s", conn.LocalAddr().String())

		backendToDial := s.Pick()

		if backendToDial == nil {
			log.Printf("no backend to dial")
			conn.Close()
			continue
		}

		log.Printf("dialing instance %s", backendToDial.GetAddr())

		s.Handle(backendToDial)

		backend, err := net.Dial("tcp", backendToDial.GetAddr())

		if err != nil {
			s.SetHealth(backendToDial, false)
			log.Println(err)
			s.Cleanup(backendToDial)
			conn.Close()
			continue
		}

		go proxy(backend, conn, func() {
			s.Snapshot()
			s.Cleanup(backendToDial)
		})
	}
}
