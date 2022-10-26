package patient

import (
	"context"
	"net/http"

	"github.com/gabriel-ross/gofhir"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Storer interface {
	Create(context.Context, gofhir.Patient) (gofhir.Patient, error)
	List(context.Context) ([]gofhir.Patient, error)
	Read(context.Context, string) (gofhir.Patient, error)
	Update(context.Context) error
	Delete(context.Context, string) error
}

type Service struct {
	router chi.Router
	storer Storer
}

func New(s Storer) *Service {
	return &Service{
		storer: s,
	}
}

func (svc *Service) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", svc.createPatient())
	r.Get("/", svc.listPatients())
	r.Route("/{patient_id}", func(r chi.Router) {
		r.Get("/", svc.getPatient())
		r.Put("/", svc.updatePatient())
		r.Delete("/", svc.deletePatient())
	})
	return r
}

func (svc *Service) createPatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		data := PatientRequest{}

		err = render.Bind(r, &data)
		if err != nil {
			render.Render(w, r, gofhir.NewError(err, http.StatusBadRequest))
			return
		}

		_, err = svc.storer.Create(ctx, data.Patient)
		if err != nil {
			render.Render(w, r, gofhir.NewError(err, http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}
}

func (svc *Service) listPatients() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error

		patients, err := svc.storer.List(ctx)
		if err != nil {
			render.Render(w, r, gofhir.NewError(err, http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.RenderList(w, r, NewPatientResponseList(patients))
		return
	}
}

func (svc *Service) getPatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error

		patient_id := chi.URLParam(r, "patient_id")
		patient, err := svc.storer.Read(ctx, patient_id)
		if status.Code(err) == codes.NotFound {
			render.Render(w, r, gofhir.NewError(err, http.StatusNotFound))
			return
		} else if err != nil {
			render.Render(w, r, gofhir.NewError(err, http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.Render(w, r, NewPatientResponse(patient))
		return
	}
}

func (svc *Service) updatePatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (svc *Service) deletePatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error

		patient_id := chi.URLParam(r, "patient_id")
		err = svc.storer.Delete(ctx, patient_id)
		if err != nil {
			render.Render(w, r, gofhir.NewError(err, http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
}
