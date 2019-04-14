package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	cursus "github.com/rodrwan/go-cursus"
	"github.com/rodrwan/go-cursus/errors"
	"github.com/rodrwan/go-cursus/response"
	"github.com/rodrwan/go-cursus/room"
	"github.com/rodrwan/go-cursus/server"
)

var upgrader websocket.Upgrader

var addr = flag.String("addr", ":8080", "service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	rooms := make(map[string]*room.Room)
	ctx := &Context{
		Rooms: rooms,
	}

	srv := server.New(*addr)
	srv.AddHandler("GET", "/ws", ws(ctx))
	srv.AddHandler("POST", "/room", createRoom(ctx))

	srv.ListenAndServe()
}

type subscription struct {
	Room string `json:"room"`
}

func createRoom(ctx *Context) server.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) (*response.Response, error) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errors.ErrInternalServerError(err)
		}
		defer r.Body.Close()

		var ss subscription
		if err := json.Unmarshal(body, &ss); err != nil {
			return nil, errors.ErrInternalServerError(err)
		}

		newRoom := room.New(ss.Room)
		go newRoom.Run()

		ctx.Rooms[ss.Room] = newRoom

		log.Printf("Room [%s] was created", ss.Room)
		return response.NewVoid(http.StatusCreated), nil
	}
}

func ws(ctx *Context) server.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) (*response.Response, error) {
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			return nil, errors.ErrInternalServerError(err)
		}

		for {
			// read incomming messages
			peerReq := &cursus.Request{}
			if err := conn.ReadJSON(peerReq); err != nil {
				return response.NewVoid(http.StatusGone), nil
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
