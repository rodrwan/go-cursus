package room

import (
	"log"

	"github.com/Finciero/cursus"
)

// Room ...
type Room struct {
	peers map[string]*cursus.Peer
	topic string

	Subscribe   chan *cursus.Peer
	Unsubscribe chan string
	Broadcast   chan *cursus.Action
}

// Run ...
func (r *Room) Run() {
	for {
		select {
		case peer := <-r.Subscribe:
			r.peers[peer.ID] = peer
		case id := <-r.Unsubscribe:
			delete(r.peers, id)
		case action := <-r.Broadcast:
			log.Println("Send message to all peers connected to users")
			log.Printf("[%s] -> %s\n", action.Type, action.Message)
			log.Printf("Message was sent to %d peers", len(r.peers))
			for _, client := range r.peers {
				client.Send(&cursus.Response{
					Message: action.Message,
				})
			}
		}
	}
}

// New ...
func New(topic string) *Room {
	return &Room{
		topic: topic,
		peers: make(map[string]*cursus.Peer),

		Subscribe:   make(chan *cursus.Peer),
		Unsubscribe: make(chan string),
		Broadcast:   make(chan *cursus.Action),
	}
}
