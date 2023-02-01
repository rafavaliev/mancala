package socket

import "mancala/mancala"

// SocketMessage stores messages sent over WS with the client who sent it
type SocketMessage struct {
	msg    []byte
	sender *Client
}

// BasePacket that can be sent between clients and server. These are all of the required attributes of any packet
type BasePacket struct {
	// Identifies the type of packet
	Type string `json:"type"`
}

// JoinPacket is sent by clients after receiving the init packet. Identifies them to the  server, and in turn other clients
type JoinPacket struct {
	BasePacket
	PlayerNumber int    `json:"player_number"` // Player number: either 0 or 1
	Slug         string `json:"slug"`
}

// Sent by clients to indicate their turn action
type TurnPacket struct {
	BasePacket

	// Player number: either 0 or 1
	PlayerNumber int `json:"player_number"`

	// room slug
	Slug string `json:"slug"`

	// Pit index
	PitIndex int `json:"pit_index"`
}

// Sent by clients to indicate their turn action
type TurnValidationPacket struct {
	BasePacket

	// Player number: either 0 or 1
	PlayerNumber int `json:"player_number"`

	// room slug
	Slug string `json:"slug"`

	// Pit index
	PitIndex int `json:"pit_index"`

	ValidationError string `json:"validation_error"`
}

// Sent by clients after confirming finishing the round
type FinishPacket struct {
	BasePacket
	// Player number: either 0 or 1
	PlayerNumber int    `json:"player_number"`
	Slug         string `json:"slug"`
}

type GetStatePacket struct {
	BasePacket
	Slug  string          `json:"slug"`
	State mancala.Mancala `json:"state"`
}
