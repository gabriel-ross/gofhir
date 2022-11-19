package bar

import "github.com/go-chi/chi"

func (svc *Service) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", svc.createBar())
	r.Get("/", svc.listBars())

	r.Route("/bar_id}", func(r chi.Router) {
		r.Get("/", svc.getBar())
		r.Put("/", svc.updateBar())
		r.Delete("/", svc.deleteBar())
	})

	return r
}
