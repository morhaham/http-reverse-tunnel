package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"sync"
)

func main() {
	tunnServAddr := "localhost:4001"
	cert, _ := tls.LoadX509KeyPair("tls/client.pem", "tls/client.key")
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", tunnServAddr, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to connect to tunneling server: %s", err)
	}
	defer conn.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	log.Printf("Connected to tunneling server on %s", tunnServAddr)
	go proxyReq(conn, &wg)
	wg.Wait()
}

func proxyReq(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	proxyTo := "localhost:4000"
	localAppConn, err := tls.Dial("tcp", proxyTo, tlsConfig)
	if err != nil {
		log.Printf("Failed to connect to HTTP server: %s", err)
		return
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err = io.Copy(localAppConn, conn)
		if err != nil {
			log.Printf("Failed to tunnel from server to local app: %s", err)
		}
	}()

	go func() {
		defer wg.Done()
		_, err = io.Copy(conn, localAppConn)
		if err != nil {
			log.Printf("Failed to tunnel from local app to server: %s", err)
		}
	}()
}
