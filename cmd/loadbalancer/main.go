package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"usman-faisal/tcp-loadbalancer/config"
)

type SafeInstanceMap struct {
	mu   sync.RWMutex
	data map[string]int
}

func initiateSafeInstanceMap(instanceList []string) *SafeInstanceMap {
	return &SafeInstanceMap{
		data: buildInstanceMap(instanceList),
	}
}

func (s *SafeInstanceMap) Increment(key string) {
	s.mu.Lock()
	s.data[key] += 1
	s.mu.Unlock()
}

func (s *SafeInstanceMap) Decrement(key string) {
	s.mu.Lock()
	s.data[key] += 1
	s.mu.Unlock()
}

func (s *SafeInstanceMap) Choose() string{
	s.mu.RLock()
	defer s.mu.RUnlock()
	var minConnectionInstance string

	for currInstance, numberConnections := range s.data {
		if minConnectionInstance == "" || numberConnections < s.data[minConnectionInstance] {
			minConnectionInstance = currInstance
		}
	}

	return minConnectionInstance
}
func (s *SafeInstanceMap) Snapshot() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap := make(map[string]int, len(s.data))
	for k, v := range s.data {
		snap[k] = v
	}
	return snap
}

func buildInstanceMap(instanceList []string) map[string]int {
	m := make(map[string]int)

	for _, val := range instanceList {
		m[val] = 0
	}

	return m
}

func (s *SafeInstanceMap) Dial(url string) (net.Conn, error) {
	conn, err := net.Dial("tcp", url)

	if err != nil {
		return nil, err
	}

	return conn, nil
}


func routeRequest(key string, instance net.Conn, conn net.Conn, safeInstanceMap *SafeInstanceMap) {
	go io.Copy(instance, conn)
	safeInstanceMap.Increment(key)

	// response
	io.Copy(conn, instance)
	safeInstanceMap.Decrement(key)

	instance.Close()
}

func main() {
	cfg, err := config.Load("config.yaml")

	if err != nil {
		log.Println(err)
		return
	}

	instanceList, port := cfg.BackendList, cfg.Port

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("listening on %s", ln.Addr().String())

	safeInstanceMap := initiateSafeInstanceMap(instanceList)

	for {
		// listen for requests
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("accepting connection %s", conn.LocalAddr().String())

		// choose which instance to dial
		instanceToDial := safeInstanceMap.Choose()

		log.Printf("dialing instance %s", instanceToDial)

		instance, err := safeInstanceMap.Dial(instanceToDial)

		if err != nil {
			log.Println(err)
		}

		// route request and increment counter
		go routeRequest(instanceToDial, instance, conn, safeInstanceMap)

		fmt.Println(safeInstanceMap.Snapshot())
	}
}
