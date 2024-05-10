package main

import (
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
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Failed to read from server: %s", err)
			return
		}
		log.Printf("Read %d bytes from server", n)
		log.Printf("Data: %s", buffer[:n])
		conn.Write([]byte("The tunnel client read the HTTP request"))
		// resp, err := http.Get(proxyTo)
		localAppConn, err := net.Dial("tcp", proxyTo)
		defer conn.Close()
		if err != nil {
			log.Printf("Failed to connect to HTTP server: %s", err)
			continue
		}
		_, err = localAppConn.Write(buffer[:n])
		if err != nil {
			log.Printf("Failed to write to HTTP server: %s", err)
			continue
		}
		// _, err = io.Copy(localAppConn, conn)
		// if err != nil {
		// 	log.Printf("Failed to read HTTP server response: %s", err)
		// 	continue
		// }
		buffer = make([]byte, 4096)
		n, err = localAppConn.Read(buffer)
		if err != nil {
			log.Printf("Failed to read HTTP server response body: %s", err)
			continue
		}
		// log.Printf("HTTP server response: %s", buffer[:n])
		conn.Write(buffer[:n])
	}
}
