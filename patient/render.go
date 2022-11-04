package patient

import (
	"net/http"

	"github.com/gabriel-ross/gofhir"
	"github.com/go-chi/render"
)

type PatientRequest struct {
	gofhir.Patient
}

func (pr *PatientRequest) Bind(r *http.Request) (err error) { return nil }

type PatientResponse struct {
	gofhir.Patient
}

func NewPatientResponse(p gofhir.Patient) *PatientResponse {
	return &PatientResponse{
		Patient: p,
	}
}

func (pr *PatientResponse) Render(w http.ResponseWriter, r *http.Request) (err error) { return nil }

func NewPatientResponseList(patients []gofhir.Patient) []render.Renderer {
	list := []render.Renderer{}
	for _, patient := range patients {
		list = append(list, NewPatientResponse(patient))
	}
	return list
}

type patientResponse struct{}

type ResponseService struct{}

func (r *ResponseService) WriteResponse() {}

func (r *ResponseService) WriteListResponse() {}
