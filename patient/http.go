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

		patients, err := svc.storer.List(ctx, *gofhir.NewListOptions())
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
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		data := PatientRequest{}

		patient_id := chi.URLParam(r, "patient_id")

		err = render.Bind(r, &data)
		if err != nil {
			render.Render(w, r, gofhir.NewError(err, http.StatusBadRequest))
			return
		}

		_, err = svc.storer.Update(ctx, patient_id, data.Patient)
		if err != nil {
			render.Render(w, r, gofhir.NewError(err, http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
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
