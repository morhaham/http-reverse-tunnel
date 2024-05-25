package main

import (
	"crypto/tls"
	"flag"
	"sync"

	proxyclient "github.com/morhaham/http-reverse-tunnel/cmd/tunnel-server/proxy-client"
)

func main() {
	listenAddr := flag.String("listenAddr", "localhost:4001", "The address to listen on for incoming connections")
	flag.Parse()
	cert, _ := tls.LoadX509KeyPair("tls/server.crt", "tls/server.key")
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	proxyClientHandler := &proxyclient.Handler{}
	proxyServer := &proxyServer{proxyClientHandler: proxyClientHandler}
	proxyServer.listen(*listenAddr, tlsConfig)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go proxyServer.accept(&wg)
	wg.Wait()
}
