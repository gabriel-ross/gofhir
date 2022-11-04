package gofhir

import (
	"errors"
	"net/http"
)

var KillRequest = errors.New("a5!C+qQMELB>C6M3;#6.")

type ErrorResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"-"`
	AppCode        int64  `json:"code,omitempty"`
	ErrorText      string `json:"error,omitempty"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.HTTPStatusCode)
	return nil
}

func NewError(err error, code int) *ErrorResponse {
	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: code,
		ErrorText:      err.Error(),
	}
}
