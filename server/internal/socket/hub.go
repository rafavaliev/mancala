package socket

import (
	"encoding/json"
	"go.uber.org/zap"
)

// Hub maintains the set of active clients
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan *SocketMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	h := &Hub{}
	h.broadcast = make(chan *SocketMessage)
	h.register = make(chan *Client)
	h.unregister = make(chan *Client)
	h.clients = map[string]*Client{}
	return h
}

// Listens for messages from websocket clients
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			zap.S().Info("A client has joined")
			h.clients[client.id] = client
		case client := <-h.unregister:
			// When a client disconnects, remove them from our clients map
			zap.S().Info("A client has disconnected")
			delete(h.clients, client.id)
			close(client.send)
		case message := <-h.broadcast:
			// Process incoming messages from clients
			zap.S().Info("Received message from client")
			h.processMessage(message)
		}
	}
}

// Sends a message to all of our clients of a certain slug
func (h *Hub) Send(slug string, msg any) {
	data, _ := json.Marshal(msg)
	h.SendBytes(slug, data)
}

// Sending bytes of data to specific room with slug
func (h *Hub) SendBytes(slug string, msg []byte) {
	for _, client := range h.clients {
		if client.slug == slug {
			client.send <- msg
		}
	}
}

// Processes an incoming message
func (h *Hub) processMessage(m *SocketMessage) {
	zap.S().With("content", string(m.msg)).Info("Processing message")
}
