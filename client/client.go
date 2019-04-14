package client

import (
	"log"
	"net/url"

	"github.com/Finciero/cursus"
	"github.com/gorilla/websocket"
)

// New ...
func New(topic string) (<-chan *cursus.Action, error) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	clientReq := &cursus.Request{
		Action: "hello",
		Topic:  topic,
	}
	if err := conn.WriteJSON(clientReq); err != nil {
		log.Println("write:", err)
		return nil, err
	}

	action := make(chan *cursus.Action)

	go func() {
		log.Println("listening message from socket")
		for {
			act := &cursus.Action{}
			err := conn.ReadJSON(act)
			if err != nil {
				log.Println("read:", err)
			}

			action <- act
		}
	}()

	return action, nil
}
