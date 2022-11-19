package bar

import (
	"context"
	"net/http"

	"github.com/gabriel-ross/gofhir"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *Service) createBar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bar create called"))
	}
}

func (svc *Service) listBars() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bar list called"))
	}
}

func (svc *Service) getBar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bar get called"))
	}
}

func (svc *Service) updateBar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bar update called"))
	}
}

func (svc *Service) deleteBar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bar delete called"))
	}
}
