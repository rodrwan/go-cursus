package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Finciero/cursus"
	"github.com/Finciero/cursus/room"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader

// Context ...
type Context struct {
	Rooms map[string]*room.Room
}

func ws(ctx *Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for {
			// read incomming messages
			clientReq := &cursus.Request{}
			if err := conn.ReadJSON(clientReq); err != nil {
				log.Println(err)
				return
			}
			// Here we need to select the corresponding topic.
			switch clientReq.Action {
			case "hello":
				clt := &cursus.Client{
					ID:     fmt.Sprintf("%d", &r),
					Socket: conn,
				}
				log.Printf("Welcome %v\n", clt)
				// insert new client into corresponding topic map.
				ctx.Rooms[clientReq.Topic].Subscribe <- clt
			case "bye":
				log.Printf("Bye client %d\n", &r)
				ctx.Rooms[clientReq.Topic].Unsubscribe <- fmt.Sprintf("%d", &r)
			case "create":
				log.Printf("create %v\n", clientReq)
				ctx.Rooms[clientReq.Topic].Broadcast <- &cursus.Action{
					Type:    clientReq.Action,
					Message: clientReq.Message,
				}
			case "update":
				log.Printf("update %v\n", clientReq)
				ctx.Rooms[clientReq.Topic].Broadcast <- &cursus.Action{
					Type:    clientReq.Action,
					Message: clientReq.Message,
				}
			case "delete":
				log.Printf("delete %v\n", clientReq)
				ctx.Rooms[clientReq.Topic].Broadcast <- &cursus.Action{
					Type:    clientReq.Action,
					Message: clientReq.Message,
				}
			}
			// send response
			clientResp := &cursus.Response{Message: "OK"}
			if err := conn.WriteJSON(clientResp); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	mux := http.NewServeMux()
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	rooms := make(map[string]*room.Room)
	userRoom := room.New("users")
	go userRoom.Run()

	orderRoom := room.New("orders")
	go orderRoom.Run()

	rooms["users"] = userRoom
	rooms["orders"] = orderRoom
	ctx := &Context{
		Rooms: rooms,
	}
	mux.HandleFunc("/ws", ws(ctx))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
