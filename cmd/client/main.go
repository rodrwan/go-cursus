package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/Finciero/cursus"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var topic = flag.String("topic", "users", "topic to subscribe")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	clientReq := &cursus.Request{
		Action:  "hello",
		Topic:   *topic,
		Message: "client",
	}
	if err := c.WriteJSON(clientReq); err != nil {
		log.Println("write:", err)
		return
	}

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			clientReq := &cursus.Request{
				Action:  "create",
				Topic:   *topic,
				Message: t.String(),
			}
			if err := c.WriteJSON(clientReq); err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			clientReq := &cursus.Request{
				Action: "bye",
				Topic:  *topic,
			}
			if err := c.WriteJSON(clientReq); err != nil {
				log.Println("write:", err)
				return
			}
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
