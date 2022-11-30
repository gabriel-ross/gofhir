package gofhir

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
)

type Server struct {
	router   chi.Router
	database *firestore.Client
	port     string
	url      string
}

func NewServer(options ...ServerOption) *Server {
	svr := &Server{}

	for _, option := range options {
		option(svr)
	}

	return svr
}

func (svr *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	svr.router.ServeHTTP(w, r)
}

type ServerOption func(*Server)
