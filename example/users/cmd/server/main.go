package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rodrwan/go-cursus"
	"github.com/rodrwan/go-cursus/emitter"
)

type Context struct {
	Emitter *emitter.Emitter
}

type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
	Status int         `json:"-"`
}

func (r *Response) Write(rw http.ResponseWriter) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(r.Status)

	return json.NewEncoder(rw).Encode(r)
}

// HandlerFunc function handler signature used by sigiriya application.
type HandlerFunc func(*Context, http.ResponseWriter, *http.Request) (*Response, error)

// Handler is an http.Handler that provides access to the Context to the given HandlerFunc.
type Handler struct {
	Ctx    *Context
	Handle HandlerFunc
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	resp, err := h.Handle(h.Ctx, rw, r)
	if err != nil {
		fmt.Println(err)
		if err = json.NewEncoder(rw).Encode(err); err != nil {
			log.Printf("[ERROR]: %v failed to encode error", err)
		}
		return
	}

	if err := resp.Write(rw); err != nil {
		log.Printf("[ERROR]: %v, encoding response: %v", err, resp)
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

	h := &Handler{
		Ctx: &Context{
			Emitter: emit,
		},
		Handle: createUser,
	}
	mux.Handle("/create", h).Methods("POST")

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

func createUser(ctx *Context, w http.ResponseWriter, r *http.Request) (*Response, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if err := ctx.Emitter.Emit(cursus.CreateAction, string(body)); err != nil {
		return nil, err
	}

	return &Response{
		Data:   "OK",
		Status: http.StatusCreated,
	}, nil
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
