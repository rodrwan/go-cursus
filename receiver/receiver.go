// Package receiver implements Subscriber interface.
package receiver

import (
	"log"
	"net/url"
	"time"

	"github.com/Finciero/cursus"
	"github.com/gorilla/websocket"
)

// Receiver represents a receiver of messages.
type Receiver struct {
	Topic string

	conn   *websocket.Conn
	action chan *cursus.Action
}

// New expose a new receiver.
func New(topic string) (*Receiver, error) {
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

	return &Receiver{
		conn:   conn,
		action: make(chan *cursus.Action),
		Topic:  topic,
	}, nil

}

// Listen listen for incomming messages.
func (r *Receiver) Listen() (<-chan *cursus.Action, error) {
	go func() {
		log.Println(">>> Listening message from socket")
		for {
			act := &cursus.Action{}
			err := r.conn.ReadJSON(act)
			if err != nil {
				close(r.action)
				return
			}

			r.action <- act
		}
	}()

	return r.action, nil
}

// Disconnect close session with server.
func (r *Receiver) Disconnect() {
	req := &cursus.Request{
		Action: "bye",
		Topic:  r.Topic,
	}
	if err := r.conn.WriteJSON(req); err != nil {
		log.Println("write:", err)
		return
	}
	err := r.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	select {
	case <-time.After(time.Second):
	}

}
