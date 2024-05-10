package main

import (
	"log"
	"net"
)

func main() {
	listenAddr := "localhost:4001"

	conn, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	log.Printf("Tunneling server started on %s", listenAddr)
	for {
		conn, err := conn.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}

		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		httpRequest := []byte("GET / HTTP/1.1\r\nHost: localhost:4000\r\n\r\n")
		_, err = conn.Write([]byte(httpRequest))
		if err != nil {
			log.Printf("Failed to write to client: %s", err)
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Failed to read from client: %s", err)
			return
		}
		log.Printf("Read %d bytes from client", n)
		log.Printf("Data: %s", buffer[:n])
	}
}
