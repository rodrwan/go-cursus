package cursus

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Request struct {
	Action  string `json:"action,omitempty"`
	Topic   string `json:"topic,omitempty"`
	Message string `json:"message,omitempty"`
}

type Response struct {
	Message string `json:"message,omitempty"`
}

type Action struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

type Client struct {
	ID     string
	Socket *websocket.Conn
	mu     sync.Mutex
}

func (c *Client) Send(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Socket.WriteJSON(v)
}
