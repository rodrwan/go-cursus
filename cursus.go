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

type Peer struct {
	ID     string
	Socket *websocket.Conn
	mu     sync.Mutex
}

func (p *Peer) Send(v interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.Socket.WriteJSON(v)
}

func (p *Peer) readPeer() {}

func (p *Peer) writePeer() {}
