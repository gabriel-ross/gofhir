package patient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabriel-ross/gofhir"
	"github.com/stretchr/testify/assert"
)

type MockPatientStorer struct{}

func (s *MockPatientStorer) Create(ctx context.Context, patient gofhir.Patient) (_ gofhir.Patient, err error) {
	return gofhir.Patient{}, nil
}

func (s *MockPatientStorer) List(ctx context.Context) (_ []gofhir.Patient, err error) {
	return []gofhir.Patient{}, nil
}

func (s *MockPatientStorer) Read(ctx context.Context, id string) (_ gofhir.Patient, err error) {
	return gofhir.Patient{}, nil
}

func (s *MockPatientStorer) Update(ctx context.Context) (err error) {
	return nil
}

func (s *MockPatientStorer) Delete(ctx context.Context, id string) (err error) {
	return nil
}

func TestCreatePatient(t *testing.T) {
	svc := New(&MockPatientStorer{})
	r := svc.Routes()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
