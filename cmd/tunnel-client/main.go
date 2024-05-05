package main

import (
	"io"
	"log"
	"net"
)

func main() {
	tunnelingSeverAddr := "tunneling-sever-ip:8080"
	conn, err := net.Dial("tcp", tunnelingSeverAddr)
	if err != nil {
		log.Fatalf("Failed to connect to tunneling server: %s", err)
	}
	defer conn.Close()

	for {
		clientConn, err := net.Dial("tcp", "localhost:3001")
		if err != nil {
			log.Printf("Failed to connect to locahost:3001: %s", err)
			continue
		}

		defer clientConn.Close()

		go func() {
			_, err := io.Copy(conn, clientConn)
			if err != nil {
				log.Printf("Error copying from clientConn to conn: %s", err)
			}
		}()

		go func() {
			_, err := io.Copy(clientConn, conn)
			if err != nil {
				log.Printf("Error copying from conn to clientConn: %s", err)
			}
		}()

	}
}
