package main

import (
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4001")
	if err != nil {
		log.Fatalf("Failed to connect to tunneling server: %s", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("Hello from the client!"))
	if err != nil {
		log.Fatalf("Failed to write to server: %s", err)
	}

	buffer := make([]byte, 100)
	_, err = conn.Read(buffer)
	if err != nil {
		log.Fatalf("Failed to read from server: %s", err)
	}
	log.Printf("Received from server: %s", buffer)

	clientAddr := conn.LocalAddr().String()

	listener, err := net.Listen("tcp", clientAddr)
	if err != nil {
		log.Fatalf("Failed to listen on client's port: %s", err)
	}
	defer listener.Close()

	log.Printf("Client listening on: %s", clientAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection: %s", err)
			continue
		}
		conn.Write([]byte("Hello from the client!"))

		go handleClient(clientConn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Failed to read from server: %s", err)
			return
		}
		log.Printf("Read %d bytes from server", n)
		log.Printf("Data: %s", buffer[:n])
	}
}
