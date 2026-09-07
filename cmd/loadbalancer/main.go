package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"usman-faisal/tcp-loadbalancer/internal/config"
	"usman-faisal/tcp-loadbalancer/internal/leastconn-balancer"
)


func initLeastConnBalancer(backendList []string)*leastconnbalancer.LeastConnBalancer {
	leastConnBalancer:=leastconnbalancer.LeastConnBalancer{}
	for _, b := range backendList {
		backend:=&leastconnbalancer.Backend{
			Addr: b,
			ActiveConns: 0,
		}

		leastConnBalancer.Heap.Push(backend)
	}
	return &leastConnBalancer
}

func proxy(backend net.Conn, conn net.Conn, lc *leastconnbalancer.LeastConnBalancer, release func()) {
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

	release()

	lc.Snapshot()
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

	lc := initLeastConnBalancer(backendList)

	for {
		// listen for requests
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("accepting connection %s", conn.LocalAddr().String())

		backendToDial:=lc.Pick()

		if backendToDial == nil {
			log.Printf("no backend to dial")
			conn.Close()
			continue
		}

		log.Printf("dialing instance %s", backendToDial.Addr)

		lc.Acquire(backendToDial)

		backend,err:=net.Dial("tcp", backendToDial.Addr)

		if err != nil {
			log.Println(err)
			lc.Release(backendToDial)
			conn.Close()
			continue
		}

		go proxy(backend, conn, lc, func() {
			lc.Release(backendToDial)
		})

	}
}
