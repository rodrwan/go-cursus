package room

import (
	"log"

	"github.com/rodrwan/go-cursus"
)

// Room reprensent a room.
// A room is where peers live, and where incomming message
// are distributed to the corresponing channel.
type Room struct {
	peers map[string]*cursus.Peer
	topic string

	Subscribe   chan *cursus.Peer
	Unsubscribe chan string
	Broadcast   chan *cursus.Action
}

// Run listen channels.
func (r *Room) Run() {
	for {
		select {
		case peer := <-r.Subscribe:
			r.peers[peer.ID] = peer
		case id := <-r.Unsubscribe:
			delete(r.peers, id)
		case action := <-r.Broadcast:
			for _, client := range r.peers {
				msg := &cursus.Response{
					Type:    action.Type,
					Message: action.Message,
				}
				client.Send(msg)
			}
			log.Printf("Message was sent to %d peers on %s", len(r.peers), r.topic)
		}
	}
}

// New returns a new Room.
func New(topic string) *Room {
	return &Room{
		topic: topic,
		peers: make(map[string]*cursus.Peer),

		Subscribe:   make(chan *cursus.Peer),
		Unsubscribe: make(chan string),
		Broadcast:   make(chan *cursus.Action),
	}
}
