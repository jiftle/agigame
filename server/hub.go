package server

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

// sendBuffer is the per-client outbound message buffer size. When full, the
// oldest frames are dropped so the emulator is never blocked.
const sendBuffer = 64

// Hub tracks websocket clients and broadcasts messages to all of them without
// ever blocking the caller.
type Hub struct {
	mu      sync.Mutex
	clients map[*Client]struct{}
}

// NewHub returns an empty Hub.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]struct{}),
	}
}

// Add registers a client and starts its write pump.
func (h *Hub) Add(conn *websocket.Conn) *Client {
	c := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, sendBuffer),
	}
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
	go c.writePump()
	return c
}

// Remove deregisters a client.
func (h *Hub) Remove(c *Client) {
	h.mu.Lock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.send)
	}
	h.mu.Unlock()
}

// BroadcastJSON marshals v and sends it to every client, dropping for any
// client whose outbound buffer is full.
func (h *Hub) BroadcastJSON(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	h.BroadcastRaw(data)
}

// BroadcastRaw sends raw JSON bytes to every client, non-blocking.
func (h *Hub) BroadcastRaw(data []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		select {
		case c.send <- data:
		default:
			// Client is slow; drop this message for it.
		}
	}
}

// Client is a single websocket connection.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// readPump drains inbound messages into the handler.
func (c *Client) readPump(handle func([]byte)) {
	defer func() {
		c.hub.Remove(c)
		_ = c.conn.Close()
	}()
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		handle(data)
	}
}

// writePump flushes the outbound queue to the socket.
func (c *Client) writePump() {
	for data := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
	_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}