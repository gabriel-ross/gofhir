package gofhir

import (
	"context"
	"net/http"
	"time"
)

type Interceptor func(context.Context, QueryEvent, func(context.Context, QueryEvent))

// Last interceptor in a chain still needs to be able to terminate -> need to wrap the calling function?

type InterceptorManager struct {
	onDatabaseTransaction []func(context.Context)
}

type QueryEvent struct {
	Time              time.Time
	Context           context.Context
	Writer            http.ResponseWriter
	Request           *http.Request
	Model             interface{}
	CollectionPath    string
	DocumentID        string
	DatabaseOperation string
	DatabaseResult    bool
	Err               error
	HTTPResponseCode  int
	HTTPResonseBody   []byte
}

// from here how do I give an interceptor access to certain parameters?

// maybe I pass in the function and it gets its parameters from the query event?

// maybe I need different methods for different types of events (like event.Database()) which would log all database events
// this would be easier if I were using a SQL database or just a database with a query language

// might need some way to translate a database function into a query

// Need a way to link events to the same query. Not the http request and response writers, but maybe some event ID

// Services take in some interface Hooker that looks like this that is a collection of sub-interfaces that it will need
type ServiceHooker interface {
	RequestHooker
}

type ServerHooker interface {
	OnServerStartup(ServerEvent)
	OnServerShutdown(ServerEvent)
}

type RequestHooker interface {
	OnRequestReceived(RequestEvent)
}

type ResponseHooker interface {
	OnServerResponse(ResponseEvent)
}

type DatabaseHooker interface {
	BeforeDatabaseQuery(DatabaseEvent)
}

type ErrorHooker interface {
	OnError(ErrorEvent)
}

type ServerEvent struct {
	Timestamp time.Time
}

type RequestEvent struct {
	RequestID string
	Timestamp time.Time
	Request   *http.Request
}

type ResponseEvent struct {
	RequestID string
	Timestamp time.Time
}

type DatabaseEvent struct {
	RequestID string
	Timestamp time.Time
}

type ErrorEvent struct {
	RequestID string
	Timestamp time.Time
}

// For now scrap endpoint-specific hooks

// hooker struct is dependency of server.
// do I need a hooker struct?
//	- I think it would just house all the registered interceptor functions

// What if I just start with implementing it on the endpoint level and seeing what code ends up being
// redundant. From there I can just abstract it out. Need to stop trying to do it perfectly the first time
// and be more okay with duplicate code and refactoring
// lets use cloud hw as demo space.
