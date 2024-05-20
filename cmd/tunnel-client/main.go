package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type AppState struct {
	apiKey string
}

func main() {

	tunnServAddr := flag.String("tunnServAddr", "localhost:4001", "The address of the tunneling server")
	apiKey := flag.String("apiKey", "1234", "The API key to use for authentication")

	app := &AppState{
		apiKey: *apiKey,
	}
	flag.Parse()
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", *tunnServAddr, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to connect to tunneling server: %s", err)
	}
	defer conn.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	log.Printf("Connected to tunneling server on %s", *tunnServAddr)
	go app.proxyReq(conn, &wg)
	wg.Wait()
}

func (app *AppState) proxyReq(conn net.Conn, wg *sync.WaitGroup) {
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
		log.Printf("api key: %s", app.apiKey)
		defer wg.Done()
		_, err = io.Copy(conn, io.MultiReader(strings.NewReader(app.apiKey+"\n"), localAppConn))
		if err != nil {
			log.Printf("Failed to tunnel from local app to server: %s", err)
		}
	}()
}
