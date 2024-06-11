package errors

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// Error is an error that formats as the given text.
type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Error returns a formatted string including the error code, error message.
func (e *Error) Error() string {
	errMsg := fmt.Sprintf("Error %d: %s", e.Code, e.Message)
	return errMsg
}

// MarshalJSON satisfies the json.Marshaler interface.
func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    e.Code,
		Message: e.Message,
	})
}

// HTTPStatus returns the http status code for the given error.
func (e *Error) HTTPStatus() int {
	return e.Code
}

// NewWrongInput returns a new WrongInput error with the given message.
func NewWrongInput(text string) *Error {
	return &Error{http.StatusBadRequest, text}
}

// NewInternalError returns a new Internal error with the given message.
func NewInternalError(text string) *Error {
	return &Error{http.StatusInternalServerError, text}
}

// NewNotFound returns a new Not Found error with the given message.
func NewNotFound(text string) *Error {
	return &Error{http.StatusNotFound, text}
}

// NewConflict returns a new Conflict error with the given message.
func NewConflict(text string) *Error {
	return &Error{http.StatusConflict, text}
}

// Encode uses the given http.ResponseWriter as a json
// encoder to response back with the appropriate http.Status
// and error body
func (e Error) Encode(ctx context.Context, w http.ResponseWriter) {
	w.WriteHeader(e.Code)
	if err := json.NewEncoder(w).Encode(&e); err != nil {
		slog.Error(fmt.Sprintf("error encoding json error response: %s", err))
	}
}
