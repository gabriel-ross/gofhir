package storage

import (
	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/gofhir"
)

type Patient struct {
	client *firestore.Client
}

func NewPatientClient(c *firestore.Client) *Patient {
	return &Patient{
		client: c,
	}
}

func (p *Patient) Create(patient gofhir.Patient) (err error) {
	return nil
}
