package socket

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"mancala/lobby"
	"mancala/mancala"
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

	// Lobby svc
	lob *lobby.Repository
	// Mancala game service
	gameService *mancala.Repository
}

func NewHub(lob *lobby.Repository, service *mancala.Repository) *Hub {
	h := &Hub{}
	h.broadcast = make(chan *SocketMessage)
	h.register = make(chan *Client)
	h.unregister = make(chan *Client)
	h.clients = map[string]*Client{}
	h.lob = lob
	h.gameService = service
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

	res := BasePacket{}
	if err := json.Unmarshal(m.msg, &res); err != nil {
		zap.S().With("error", err).Error("received invalid json message in processMessage!")
		return
	}
	ctx := context.Background()

	switch res.Type {
	case "join":
		// each player sends a join packet when joined to populate room
		// game begins with the second player joins the room
		log.Println("Received join packet")
		res := JoinPacket{}
		json.Unmarshal(m.msg, &res)
		h.clients[m.sender.id].slug = res.Slug

		// if the player is the first to join, they are player #0, otherwise #1
		count := 0
		for _, client := range h.clients {
			if client.slug == res.Slug {
				count++
			}
		}
		h.clients[m.sender.id].slug = res.Slug
		if count == 0 {
			m.sender.playerNumber = 0
		} else {
			m.sender.playerNumber = 1
		}

		h.Send(res.Slug, res)

	case "state":
		// represents the state of the game
		log.Println("Received state packet")
		res := GetStatePacket{}
		json.Unmarshal(m.msg, &res)

		m, err := h.gameService.Get(ctx, res.Slug)
		if err != nil {
			m = mancala.Start(res.Slug)
		}
		h.Send(res.Slug, m)
	case "turn":
		// represents player's turn
		log.Println("Received turn packet")
		res := TurnPacket{}
		json.Unmarshal(m.msg, &res)

		m, err := h.gameService.Get(ctx, res.Slug)
		if err != nil {
			m = mancala.Start(res.Slug)
		}
		err = m.PlayTurn(mancala.PlayerNumber(res.PlayerNumber), res.PitIndex)
		if err != nil {
			// send validation error back to client
			return
		}
		zap.S().With("board", m.Board).Info("new state")

		// save game state
		_, _ = h.gameService.Save(ctx, m)
		state := GetStatePacket{
			BasePacket: BasePacket{Type: "state"},
			Slug:       res.Slug,
			State:      *m,
		}
		h.Send(res.Slug, state)

	// confirmation of round completion
	case "finish":
		log.Println("Received finish packet")
		res := FinishPacket{}
		json.Unmarshal(m.msg, &res)

		h.Send(res.Slug, res)
	}
}
