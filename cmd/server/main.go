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
			peerReq := &cursus.Request{}
			if err := conn.ReadJSON(peerReq); err != nil {
				log.Println(err)
				return
			}
			// Here we need to select the corresponding topic.
			switch peerReq.Action {
			case "hello":
				peer := &cursus.Peer{
					ID:     fmt.Sprintf("%d", &r),
					Socket: conn,
				}
				log.Printf("Welcome %s\n", peer.ID)
				// insert new peer into corresponding topic map.
				ctx.Rooms[peerReq.Topic].Subscribe <- peer
			case "bye":
				log.Printf("Bye %d\n", &r)
				ctx.Rooms[peerReq.Topic].Unsubscribe <- fmt.Sprintf("%d", &r)
			case "create":
				log.Printf("create %v\n", peerReq)
				ctx.Rooms[peerReq.Topic].Broadcast <- &cursus.Action{
					Type:    peerReq.Action,
					Message: peerReq.Message,
				}
				continue
			case "update":
				log.Printf("update %v\n", peerReq)
				ctx.Rooms[peerReq.Topic].Broadcast <- &cursus.Action{
					Type:    peerReq.Action,
					Message: peerReq.Message,
				}
				continue
			case "delete":
				log.Printf("delete %v\n", peerReq)
				ctx.Rooms[peerReq.Topic].Broadcast <- &cursus.Action{
					Type:    peerReq.Action,
					Message: peerReq.Message,
				}
				continue
			}
			// send response
			peerResp := &cursus.Response{Message: "OK"}
			if err := conn.WriteJSON(peerResp); err != nil {
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
