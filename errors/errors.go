package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse represents a server error.
type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func (er *ErrorResponse) Write(rw http.ResponseWriter) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(er.Status)

	return json.NewEncoder(rw).Encode(er)
}

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("[Message]: %s | [Status]: %d", er.Message, er.Status)
}

// New returns a new error.
func New(message string, status int) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Status:  status,
	}
}

// ErrInternalServerError ...
func ErrInternalServerError(err error) *ErrorResponse {
	return &ErrorResponse{
		Message: err.Error(),
		Status:  http.StatusInternalServerError,
	}
}
