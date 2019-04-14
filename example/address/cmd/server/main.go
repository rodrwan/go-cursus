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
	"syscall"
	"time"

	"github.com/Finciero/cursus/emitter"
	"github.com/Finciero/cursus/example/users"
	"github.com/Finciero/cursus/receiver"
)

type Context struct {
	Emitter *emitter.Emitter
}

func createAddress(ctx *Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user users.Address

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}

		if err := json.Unmarshal(body, &user); err != nil {
			return
		}

		w.WriteHeader(201)
	}
}

var addr = flag.String("addr", ":8082", "service address")

func main() {
	log.SetFlags(0)

	mux := http.NewServeMux()

	receiver, err := receiver.New("users")
	if err != nil {
		log.Fatal(err)
	}
	defer receiver.Disconnect()

	action, err := receiver.Listen()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case act, close := <-action:
				if !close {
					return
				}

				if act.Type == "create" {
					var u users.User
					if err := json.Unmarshal([]byte(act.Message), &u); err != nil {
						continue
					}
					log.Printf("Address: %s\n", u.Address.Street)
				}
			}
		}
	}()

	mux.HandleFunc("/create", createAddress(&Context{}))

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
