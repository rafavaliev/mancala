package socket

// SocketMessage stores messages sent over WS with the client who sent it
type SocketMessage struct {
	msg    []byte
	sender *Client
}
