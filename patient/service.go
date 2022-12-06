package patient

import (
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/gofhir/hook"
	"github.com/go-chi/chi"
)

type Service struct {
	router         chi.Router
	db             *firestore.Client
	absolutePath   string
	requestCounter int

	RequestInterceptors  []hook.RequestInterceptor
	ResponseInterceptors []hook.ResponseInterceptor
	DatabaseInterceptors []hook.DatabaseInterceptor
}

// New returns a new User service and mounts it at "/slug". slug should
// not contain a leading slash and baseURL should not contain a trailing
// slash.
func New(router chi.Router, db *firestore.Client, baseURL, slug string) *Service {
	svc := &Service{
		router:         router,
		db:             db,
		absolutePath:   fmt.Sprintf("%s/%s", baseURL, slug),
		requestCounter: 0,
	}
	router.Mount("/"+slug, svc.Routes())
	return svc
}

func (svc *Service) NewRequestID() int {
	svc.requestCounter++
	return svc.requestCounter
}

func (svc *Service) RegisterInterceptors(interceptors ...interface{}) {
	for _, i := range interceptors {
		if interceptor, ok := i.(hook.RequestInterceptor); ok {
			log.Println("Successfully registered request interceptor")
			svc.RequestInterceptors = append(svc.RequestInterceptors, interceptor)
		}
		if interceptor, ok := i.(hook.ResponseInterceptor); ok {
			log.Println("Successfully registered response interceptor")
			svc.ResponseInterceptors = append(svc.ResponseInterceptors, interceptor)
		}
		if interceptor, ok := i.(hook.DatabaseInterceptor); ok {
			log.Println("Successfully registered database interceptor")
			svc.DatabaseInterceptors = append(svc.DatabaseInterceptors, interceptor)
		}
	}
}
