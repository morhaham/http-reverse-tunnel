package main

import (
	"log"
	"net"
)

func main() {
	listenAddr := "localhost:4001"

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	defer listener.Close()

	log.Printf("Tunneling server started on %s", listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}
		readHello := make([]byte, 100)
		_, err = conn.Read(readHello)
		if err != nil {
			log.Printf("Failed to read from client: %s", err)
		}
		log.Printf("Received from client: %s", readHello)
		_, err = conn.Write([]byte("Hello from the server!"))
		if err != nil {
			log.Printf("Failed to write to client: %s", err)
		}
		// targetAddr := conn.RemoteAddr().String()
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

