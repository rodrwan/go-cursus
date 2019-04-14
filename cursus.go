package cursus

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Allowed actions
const (
	CreateAction = "create"
	UpdateAction = "update"
	DeleteAction = "delete"
)

// Request ...
type Request struct {
	Action  string `json:"action,omitempty"`
	Topic   string `json:"topic,omitempty"`
	Message string `json:"message,omitempty"`
}

// Response ...
type Response struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

// Action ...
type Action struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

// Peer ...
type Peer struct {
	ID     string
	Socket *websocket.Conn
	mu     sync.Mutex
}

// Send send message as json on socket. This method has mutex to avoid channel saturation.
func (p *Peer) Send(v interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.Socket.WriteJSON(v)
}

// Publisher ...
type Publisher interface {
	Emit(string) error
	Disconnect()
}

// Subscriber ...
type Subscriber interface {
	Listen() (<-chan *Action, error)
	Disconnect()
}
