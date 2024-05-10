package main

import (
	"io"
	"log"
	"net"
	"net/http"
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

	proxyTo := "http://localhost:4000"
	for {
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Failed to read from server: %s", err)
			return
		}
		log.Printf("Read %d bytes from server", n)
		log.Printf("Data: %s", buffer[:n])
		conn.Write([]byte("The tunnel client read the HTTP request"))
		resp, err := http.Get(proxyTo)
		if err != nil {
			log.Printf("Failed to read HTTP server response: %s", err)
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read HTTP server response body: %s", err)
			continue
		}
		log.Printf("HTTP server response: %s", respBody)
	}
}
