package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"usman-faisal/tcp-loadbalancer/internal/config"
	"usman-faisal/tcp-loadbalancer/internal/scheduler"
)

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

		go func(conn net.Conn) {
			backendToDial, err := s.Submit(conn)

			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}

			log.Printf("dialing instance %s", backendToDial.GetAddr())
		}(conn)

	}
}
