package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rodrwan/go-cursus/errors"
	"github.com/rodrwan/go-cursus/httprouter"
	"github.com/rodrwan/go-cursus/response"
	"github.com/urfave/negroni"
)

// HandlerFunc function handler signature used by sigiriya application.
type HandlerFunc func(http.ResponseWriter, *http.Request) (*response.Response, error)

// Handler represents a handler that join context with handlerfunc.
type Handler struct {
	Handle HandlerFunc
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	resp, err := h.Handle(rw, r)
	if err != nil {
		if e, ok := (err).(*errors.ErrorResponse); ok {
			if err = e.Write(rw); err != nil {
				log.Printf("[ERROR]: %v, encoding error: %v", err, e)
			}
			return
		}

		if err = json.NewEncoder(rw).Encode(err); err != nil {
			log.Printf("[ERROR]: %v failed to encode error", err)
		}
		return
	}

	if resp.Status != http.StatusGone {
		if err := resp.Write(rw); err != nil {
			log.Printf("[ERROR]: %v, encoding response: %v", err, resp)
		}
	}
}

// Server ...
type Server struct {
	Router *httprouter.Router
	Server *http.Server
}

// AddHandler ...
func (s *Server) AddHandler(method, path string, hf HandlerFunc) {
	h := &Handler{
		Handle: hf,
	}
	switch method {
	case "GET":
		s.Router.GET(path, h)
	case "POST":
		s.Router.POST(path, h)
	}
}

// ListenAndServe ...
func (s *Server) ListenAndServe() {

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())

	n.UseHandler(s.Router)
	s.Server.Handler = n

	log.Printf("Listening on: %s", s.Server.Addr)
	log.Fatal(s.Server.ListenAndServe())
}

// New ...
func New(addr string) *Server {
	router := httprouter.New()

	server := &http.Server{
		Addr: addr,
	}

	return &Server{
		Router: router,
		Server: server,
	}
}
