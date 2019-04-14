package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Finciero/cursus"
	"github.com/Finciero/cursus/room"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader

// Context ...
type Context struct {
	Rooms map[string]*room.Room
}

type subscription struct {
	Room string `json:"room"`
}

func createRoom(ctx *Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		defer r.Body.Close()

		var ss subscription
		if err := json.Unmarshal(body, &ss); err != nil {
			return
		}

		newRoom := room.New(ss.Room)
		go newRoom.Run()

		ctx.Rooms[ss.Room] = newRoom

		w.WriteHeader(http.StatusCreated)
	}
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
				ctx.Rooms[peerReq.Topic].Broadcast <- &cursus.Action{
					Type:    peerReq.Action,
					Message: peerReq.Message,
				}
			case "update":
				ctx.Rooms[peerReq.Topic].Broadcast <- &cursus.Action{
					Type:    peerReq.Action,
					Message: peerReq.Message,
				}
			case "delete":
				ctx.Rooms[peerReq.Topic].Broadcast <- &cursus.Action{
					Type:    peerReq.Action,
					Message: peerReq.Message,
				}
			}
		}
	}
}

func main() {
	log.SetFlags(0)

	mux := mux.NewRouter()
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

	mux.HandleFunc("/room", createRoom(ctx)).Methods("POST")

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
