package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"time"

	"github.com/Finciero/cursus/emitter"
	"github.com/gorilla/mux"
)

type Context struct {
	Emitter *emitter.Emitter
}

func createUser(ctx *Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user user.User

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(body, &user); err != nil {
			return
		}

		ctx.Emitter.Emit("create", string(body))

		w.WriteHeader(http.StatusCreated)
	}
}

var addr = flag.String("addr", ":8081", "service address")

func main() {
	log.SetFlags(0)

	mux := mux.NewRouter()

	emit, err := emitter.New("users")
	if err != nil {
		log.Fatal(err)
	}

	ctx := &Context{
		Emitter: emit,
	}
	mux.HandleFunc("/create", createUser(ctx)).Methods("POST")

	server := http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		log.Printf("Listening on http://0.0.0.0%s\n", *addr)
		log.Fatal(server.ListenAndServe())
	}()

	<-interrupt
	emit.Disconnect()

	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx1)
}
