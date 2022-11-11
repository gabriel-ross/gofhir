package gofhir

import (
	"context"
	"net/http"
	"time"
)

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

// Do I really need some way to link a database transaction to an http request? It's consuming a lot of brainpower
