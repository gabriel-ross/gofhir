package patient

import (
	"context"
	"net/http"

	"github.com/gabriel-ross/gofhir"
	"github.com/go-chi/chi"
)

type Storer interface {
	Create(context.Context, gofhir.Patient) (gofhir.Patient, error)
	List(context.Context, gofhir.ListOptions) ([]gofhir.Patient, error)
	Read(context.Context, string) (gofhir.Patient, error)
	Update(context.Context, string, gofhir.Patient) (gofhir.Patient, error)
	Delete(context.Context, string) error
}

type Renderer interface {
	Render(http.ResponseWriter, *http.Request, gofhir.Patient, string, int, map[string]string) error
	RenderList(http.ResponseWriter, *http.Request, []gofhir.Patient, string, int, map[string]string) error
}

type Service struct {
	router   chi.Router
	path     string
	storer   Storer
	renderer Renderer
}

func New(storer Storer, renderer Renderer, path string, r chi.Router) *Service {
	svc := &Service{
		path:   path,
		storer: storer,
	}
	svc.router = svc.Routes()
	r.Mount(path, svc.router)
	return svc
}
