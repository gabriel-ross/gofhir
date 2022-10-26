package gofhir

import (
	"net/http"
)

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
