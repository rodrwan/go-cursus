// Package emitter implements Publisher interface.
package emitter

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rodrwan/go-cursus"
)

// Emitter represents a emitter of messages.
type Emitter struct {
	conn  *websocket.Conn
	Topic string
}

// New expose a new Emitter.
func New(topic string) (*Emitter, error) {
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

	return &Emitter{
		conn:  conn,
		Topic: topic,
	}, nil
}

// Emit emit messages to peers.
func (e *Emitter) Emit(action, message string) error {
	clientReq := &cursus.Request{
		Action:  action,
		Topic:   e.Topic,
		Message: message,
	}
	if err := e.conn.WriteJSON(clientReq); err != nil {
		return err
	}

	return nil
}

// Disconnect close session with server.
func (e *Emitter) Disconnect() {
	req := &cursus.Request{
		Action: "bye",
		Topic:  e.Topic,
	}
	if err := e.conn.WriteJSON(req); err != nil {
		log.Println("write:", err)
		return
	}
	err := e.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	select {
	case <-time.After(time.Second):
	}
}
