package websocket

import (
	"log"
	"sync"
)

// Hub coordinates real-time event distribution safely across multi-threaded routines
type Hub struct {
	clients    map[string]*Client
	mu         sync.RWMutex
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

// Run executes a blocking event-loop listener handling memory allocations concurrently
func (h *Hub) Run() {
	log.Println("[WEBSOCKET] Thread-safe network orchestration hub active.")
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("[WEBSOCKET] Client linked cleanly. User Session: %s", client.UserID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, exists := h.clients[client.UserID]; exists {
				delete(h.clients, client.UserID)
				close(client.Send)
				log.Printf("[WEBSOCKET] Client disconnected safely. Evicted User: %s", client.UserID)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.UserID)
				}
			}
			h.mu.RUnlock()
		}
	}
}