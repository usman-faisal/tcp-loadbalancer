package main

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var success, failed int64
	var wg sync.WaitGroup
	concurrency := 20000
	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := net.DialTimeout("tcp", "localhost:8000", 3*time.Second)
			if err != nil {
				atomic.AddInt64(&failed, 1)
				return
			}
			defer conn.Close()
			conn.Write([]byte("ping\n"))
			buf := make([]byte, 64)
			conn.SetReadDeadline(time.Now().Add(3 * time.Second))
			_, err = conn.Read(buf)
			if err != nil {
				atomic.AddInt64(&failed, 1)
				return
			}
			atomic.AddInt64(&success, 1)
		}()
	}

	wg.Wait()
	fmt.Printf("done in %v: %d success, %d failed\n", time.Since(start), success, failed)
}
