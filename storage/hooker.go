package storage

import (
	"context"

	"github.com/gabriel-ross/gofhir"
)

type Hooker struct {
	svc                 Service
	beforeDBTransaction []func()
}

func NewHooker(s Service) *Hooker {
	return &Hooker{
		svc: s,
	}
}

func (h *Hooker) RegisterBeforeDBTransactionInterceptor(fn func()) {
	h.beforeDBTransaction = append(h.beforeDBTransaction, fn)
}

func (h *Hooker) Create(ctx context.Context, patient gofhir.Patient) (_ gofhir.Patient, err error) {
	for _, fn := range h.beforeDBTransaction {
		fn()
	}
	return h.svc.Create(ctx, patient)
}

func (h *Hooker) List(ctx context.Context) (_ []gofhir.Patient, err error) {
	for _, fn := range h.beforeDBTransaction {
		fn()
	}
	return h.svc.List(ctx)
}

func (h *Hooker) Read(ctx context.Context, id string) (_ gofhir.Patient, err error) {
	for _, fn := range h.beforeDBTransaction {
		fn()
	}
	return h.svc.Read(ctx, id)
}

func (h *Hooker) Update(ctx context.Context) (err error) {
	for _, fn := range h.beforeDBTransaction {
		fn()
	}
	return h.svc.Update(ctx)
}

func (h *Hooker) Delete(ctx context.Context, id string) (err error) {
	for _, fn := range h.beforeDBTransaction {
		fn()
	}
	return h.svc.Delete(ctx, id)
}
