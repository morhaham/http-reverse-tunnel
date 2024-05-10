package main

import (
	"io"
	"log"
	"net"
	"sync"
)

func main() {
	tunnServAddr := "localhost:4001"
	conn, err := net.Dial("tcp", tunnServAddr)
	if err != nil {
		log.Fatalf("Failed to connect to tunneling server: %s", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	log.Printf("Connected to tunneling server on %s", tunnServAddr)
	go proxyReq(conn)
	wg.Wait()
}

func proxyReq(conn net.Conn) {
	defer conn.Close()

	proxyTo := "localhost:4000"
	for {
		localAppConn, err := net.Dial("tcp", proxyTo)
		defer conn.Close()
		if err != nil {
			log.Printf("Failed to connect to HTTP server: %s", err)
			continue
		}
		go func() {
			_, err = io.Copy(localAppConn, conn)
			if err != nil {
				log.Printf("Failed to tunnel from server to local app: %s", err)
				return
			}
		}()

		_, err = io.Copy(conn, localAppConn)
		if err != nil {
			log.Printf("Failed to tunnel from local app to server: %s", err)
			return
		}
	}
}
