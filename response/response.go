package response

import (
	"encoding/json"
	"net/http"
)

// Response represent a standarized server response .
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

// New returns a Response with data.
func New(data interface{}, status int) *Response {
	return &Response{
		Data:   data,
		Status: status,
	}
}

// NewWithMeta returns a Response with data and meta.
func NewWithMeta(data, meta interface{}, status int) *Response {
	return &Response{
		Data:   data,
		Meta:   meta,
		Status: status,
	}
}

// NewVoid returns a Response only with status.
func NewVoid(status int) *Response {
	return &Response{
		Status: status,
	}
}
