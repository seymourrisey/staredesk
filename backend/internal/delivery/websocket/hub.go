package websocketdelivery

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client merepresentasikan satu koneksi WebSocket aktif.
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Hub mengelola semua client yang terkoneksi.
// Semua mutasi state dilakukan melalui channel — goroutine-safe.
type Hub struct {
	clients    map[*Client]struct{}
	mu         sync.RWMutex
	register   chan *Client
	unregister chan *Client
	Broadcast  chan []byte
}

// NewHub membuat Hub baru siap pakai.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Broadcast:  make(chan []byte, 256),
	}
}

// Run harus dijalankan dalam goroutine tersendiri.
// Mengelola register, unregister, dan broadcast secara serial.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = struct{}{}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Channel penuh — client lambat, drop dan unregister
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}
