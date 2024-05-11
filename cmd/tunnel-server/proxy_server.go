package main

import (
	"crypto/tls"
	"log"
	"net"
	"sync"

	proxyclient "github.com/morhaham/http-reverse-tunnel/cmd/tunnel-server/proxy-client"
)

type proxyServer struct {
	listener           net.Listener
	proxyClientHandler *proxyclient.Handler
	// httpClientHandler httpclient.Handler
}

func (ps *proxyServer) listen(listenAddr string, tlsConfig *tls.Config) {
	listener, err := tls.Listen("tcp", listenAddr, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	log.Printf("Tunneling server started on %s", listenAddr)
	ps.listener = listener
}

func (ps *proxyServer) accept(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		proxyClientConn, err := ps.listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}
		log.Printf("Accepted connection from %s", proxyClientConn.RemoteAddr())
		if err != nil {
			log.Printf("Failed to proxy HTTP request: %s", err)
			continue
		}
		ps.proxyClientHandler.Conn = proxyClientConn
		go ps.proxyClientHandler.HandleConn()
	}
}
