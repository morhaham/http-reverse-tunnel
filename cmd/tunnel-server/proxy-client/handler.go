package proxyclient

import (
	"log"
	"net"
	"sync"
)

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

func (h *Handler) authenticate() {

}

// TODO: proxy traffic from external http clients
func (h *Handler) handleWrites(wg *sync.WaitGroup) {
	defer wg.Done()
	httpRequest := []byte("GET / HTTP/1.1\r\nHost: localhost:4000\r\n\r\n")
	_, err := h.Conn.Write([]byte(httpRequest))
	if err != nil {
		log.Printf("Failed to write to client: %s", err)
	}
}

func (h *Handler) handleReads(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		buffer := make([]byte, 4096)
		n, err := h.Conn.Read(buffer)
		if err != nil {
			log.Printf("Failed to read from client: %s", err)
			return
		}
		log.Printf("Read %d bytes from client", n)
		log.Printf("Data: %s", buffer[:n])
	}
}
