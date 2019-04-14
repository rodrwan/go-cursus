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

func New(data interface{}, status int) *Response {
	return &Response{
		Data:   data,
		Status: status,
	}
}

func NewWithMeta(data, meta interface{}, status int) *Response {
	return &Response{
		Data:   data,
		Meta:   meta,
		Status: status,
	}
}

func NewVoid(status int) *Response {
	return &Response{
		Status: status,
	}
}
