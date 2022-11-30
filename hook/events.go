package hook

import (
	"context"
	"net/http"
	"time"
)

// NEED THESE IF I UNEXPORT Context
//
// func NewRequestEvent(e *Context, r *http.Request, w http.ResponseWriter) *RequestEvent {
// 	return &RequestEvent{
// 		Context:      e,
// 		Request:        r,
// 		ResponseWriter: w,
// 	}
// }

// func NewResponseEvent(e *Context, r *http.Request, w http.ResponseWriter, code int, responseBody []byte) *ResponseEvent {
// 	return &ResponseEvent{
// 		Context:       e,
// 		Request:         r,
// 		ResponseWriter:  w,
// 		HTTPStatusCode:  code,
// 		HTTPResonseBody: responseBody,
// 	}
// }

type Context struct {
	RequestID   string
	ShouldAbort bool
}

func NewContext(id string) *Context {
	return &Context{
		RequestID:   id,
		ShouldAbort: false,
	}
}

func (e *Context) Abort() {
	e.ShouldAbort = true
}

type ServerEvent struct {
	Timestamp time.Time
}

type RequestEvent struct {
	*Context
	Timestamp      time.Time
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

type ResponseEvent struct {
	*Context
	Timestamp       time.Time
	Request         *http.Request
	ResponseWriter  http.ResponseWriter
	HTTPStatusCode  int
	HTTPResonseBody []byte
}

type DatabaseQueryEvent struct {
	*Context
	Timestamp  time.Time
	Ctx        context.Context
	Query      string
	Successful bool
	Result     interface{}
	Error      error
}

type ErrorEvent struct {
	*Context
	Timestamp        time.Time
	Ctx              context.Context
	Request          *http.Request
	ResponseWriter   http.ResponseWriter
	DatabaseQuery    string
	DatabaseResult   map[string]interface{}
	Err              error
	HTTPResponseCode int
	HTTPResonseBody  []byte
}

// type QueryEvent struct {
// 	Time              time.Time
// 	Context           context.Context
// 	Writer            http.ResponseWriter
// 	Request           *http.Request
// 	Model             interface{}
// 	CollectionPath    string
// 	DocumentID        string
// 	DatabaseOperation string
// 	DatabaseResult    bool
// 	Err               error
// 	HTTPResponseCode  int
// 	HTTPResonseBody   []byte
// }
