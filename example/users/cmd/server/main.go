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
	"syscall"
	"time"

	"github.com/rodrwan/go-cursus/emitter"
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

		if err := ctx.Emitter.Emit("create", string(body)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
	defer emit.Disconnect()

	ctx := &Context{
		Emitter: emit,
	}
	mux.HandleFunc("/create", createUser(ctx)).Methods("POST")

	server := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	go func() {
		log.Printf("Listening on http://0.0.0.0%s\n", *addr)
		log.Fatal(server.ListenAndServe())
	}()

	graceful(server)
}

func graceful(hs *http.Server) {
	stop := make(chan os.Signal, 1)
	timeout := 5 * time.Second

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("\nShutdown with timeout: %s\n", timeout)

	if err := hs.Shutdown(ctx); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Println("Server stopped")
	}
}
