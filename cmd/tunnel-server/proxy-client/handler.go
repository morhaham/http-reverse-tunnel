package proxyclient

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

var AuthenticationError = errors.New("Failed to authenticate client")

type Handler struct {
	Conn net.Conn
}

func (h *Handler) HandleConn() {
	defer h.Conn.Close()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go h.handleReads(&wg)
	go h.handleWrites(&wg)
	wg.Wait()
}

func (h *Handler) authenticate(apiKey string) error {
	expectedApiKey := "1234"
	if apiKey != expectedApiKey {
		log.Printf("apiKey: '%s' (len: %d)", apiKey, len(apiKey))
		log.Printf("expectedApiKey: '%s' (len: %d)", expectedApiKey, len(expectedApiKey))
		return AuthenticationError
	}
	return nil
}

func (h *Handler) handleWrites(wg *sync.WaitGroup) {
	defer wg.Done()
	httpRequest := []byte("GET / HTTP/1.1\r\nHost: localhost:4000\r\n\r\n")
	_, err := h.Conn.Write(httpRequest)
	if err != nil {
		log.Printf("Failed to write to client: %s", err)
	}
}

func (h *Handler) handleReads(wg *sync.WaitGroup) {
	defer wg.Done()
	reader := bufio.NewReader(h.Conn)

	for {
		// Read and authenticate the API key for each request
		apiKey, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading API key: %s", err)
			return
		}

		trimmedApiKey := strings.TrimSpace(apiKey)
		err = h.authenticate(trimmedApiKey)
		if err != nil {
			log.Print(err)
			h.Conn.Close()
			return
		}

		// Read the HTTP headers to get the Content-Length
		headers := make(map[string]string)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading header: %s", err)
				return
			}
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		// Get the Content-Length and read the body accordingly
		contentLength, err := strconv.Atoi(headers["Content-Length"])
		if err != nil {
			log.Printf("Error converting Content-Length: %s", err)
			return
		}

		body := make([]byte, contentLength)
		n, err := io.ReadFull(reader, body)
		if err != nil {
			log.Printf("Error reading body: %s", err)
			return
		}
		log.Printf("Read %d bytes", n)
		log.Printf("Received full response: %s", string(body))
	}
}
