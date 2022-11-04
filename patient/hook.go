package patient

import (
	"context"
	"net/http"

	"github.com/gabriel-ross/gofhir"
)

type Hooker struct {
	onHTTPRequest         []func()
	onGetRequest          []func()
	onDatabaseTransaction []func()
	onDatabaseCreate      []func()

	storer   Storer
	renderer Renderer
}

func (h *Hooker) CatchHTTPRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (h *Hooker) Create(ctx context.Context, patient gofhir.Patient) (_ gofhir.Patient, err error) {
	return h.storer.Create(ctx, patient)
}

func (h *Hooker) List(ctx context.Context, opts gofhir.ListOptions) (_ []gofhir.Patient, err error) {
	return h.storer.List(ctx, opts)
}

func (h *Hooker) Read(ctx context.Context, id string) (_ gofhir.Patient, err error) {
	return h.storer.Read(ctx, id)
}

func (h *Hooker) Update(ctx context.Context, id string, patient gofhir.Patient) (_ gofhir.Patient, err error) {
	return h.storer.Update(ctx, id, patient)
}

func (h *Hooker) Delete(ctx context.Context, id string) (err error) {
	return h.storer.Delete(ctx, id)
}

func (h *Hooker) Render(w http.ResponseWriter, r *http.Request, patient gofhir.Patient, path string, code int, headers map[string]string) (err error) {
	return nil
}

func (h *Hooker) RenderList(w http.ResponseWriter, r *http.Request, patients []gofhir.Patient, path string, code int, headers map[string]string) (err error) {
	return nil
}
