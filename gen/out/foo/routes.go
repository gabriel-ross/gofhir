package foo

import "github.com/go-chi/chi"

func (svc *Service) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", svc.createFoo())
	r.Get("/", svc.listFoos())

	r.Route("/foo_id}", func(r chi.Router) {
		r.Get("/", svc.getFoo())
		r.Put("/", svc.updateFoo())
		r.Delete("/", svc.deleteFoo())
	})

	return r
}
