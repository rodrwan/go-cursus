package room

import (
	"log"

	"github.com/Finciero/cursus"
)

// Room ...
type Room struct {
	clients map[string]*cursus.Client
	topic   string

	Subscribe   chan *cursus.Client
	Unsubscribe chan string
	Broadcast   chan *cursus.Action
}

// Run ...
func (r *Room) Run() {
	for {
		select {
		case client := <-r.Subscribe:
			r.clients[client.ID] = client
		case id := <-r.Unsubscribe:
			delete(r.clients, id)
		case action := <-r.Broadcast:
			log.Println("Send message to all peers connected to users")
			log.Printf("[%s] -> %s\n", action.Type, action.Message)
			log.Printf("Message was sent to %d peers", len(r.clients))
			for _, client := range r.clients {
				client.Send(&cursus.Response{
					Message: action.Message,
				})
			}
		}
	}
}

func New(topic string) *Room {
	return &Room{
		topic:   topic,
		clients: make(map[string]*cursus.Client),

		Subscribe:   make(chan *cursus.Client),
		Unsubscribe: make(chan string),
		Broadcast:   make(chan *cursus.Action),
	}
}
