package patient

import (
	"net/http"
)

func (svc *Service) createPatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("patient create called"))
	}
}

func (svc *Service) listPatients() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("patient list called"))
	}
}

func (svc *Service) getPatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("patient get called"))
	}
}

func (svc *Service) updatePatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("patient update called"))
	}
}

func (svc *Service) deletePatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("patient delete called"))
	}
}
