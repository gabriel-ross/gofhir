package patient

import "github.com/go-chi/chi"

func (svc *Service) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use()

	r.Post("/", svc.createPatient())
	r.Get("/", svc.listPatients())

	r.Route("/{patient_id}", func(r chi.Router) {
		r.Get("/", svc.getPatient())
		r.Put("/", svc.updatePatient())
		r.Delete("/", svc.deletePatient())
	})

	return r
}
