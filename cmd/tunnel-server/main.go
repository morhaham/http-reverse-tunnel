package main

import (
	"io"
	"log"
	"net"
)

func handleClient(conn net.Conn, targetAddr string) {
	defer conn.Close()

	// Connect to the target address
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Failed to connect to target address: %s", err)
		return
	}
	defer targetConn.Close()

	// Start bi-directional data transfer between client and target
	go func() {
		_, err := io.Copy(conn, targetConn)
		if err != nil {
			log.Printf("Error copying from target to client: %s", err)
		}
	}()

	_, err = io.Copy(targetConn, conn)
	if err != nil {
		log.Printf("Error copying from client to target: %s", err)
	}
}

func main() {
	// Define the listening address and port
	listenAddr := "localhost:8080"
	targetAddr := "localhost:3001"

	// Start listening for incoming connections
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	defer listener.Close()

	log.Printf("Tunneling server started on %s, forwarding traffic to %s", listenAddr, targetAddr)

	// Accept and handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}

		// Handle incoming connection in a separate goroutine
		go handleClient(conn, targetAddr)
	}
}
