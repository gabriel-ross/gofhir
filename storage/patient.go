package storage

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/gofhir"
	"google.golang.org/api/iterator"
)

type Service struct {
	client     *firestore.Client
	collection string
}

func NewPatientClient(c *firestore.Client) *Service {
	return &Service{
		client:     c,
		collection: "patient",
	}
}

func (s *Service) Create(ctx context.Context, patient gofhir.Patient) (_ gofhir.Patient, err error) {
	ref := s.client.Collection(s.collection).NewDoc()
	patient.ID = ref.ID
	_, err = s.client.Collection(s.collection).Doc(ref.ID).Set(ctx, patient)
	return patient, err
}

func (s *Service) List(ctx context.Context) (_ []gofhir.Patient, err error) {
	resp := []gofhir.Patient{}
	iter := s.client.Collection(s.collection).Documents(ctx)
	for {
		dsnap, done := iter.Next()
		if done == iterator.Done {
			break
		}
		var patient gofhir.Patient
		dsnap.DataTo(&patient)
		resp = append(resp, patient)
	}
	return resp, nil
}

func (s *Service) Read(ctx context.Context, id string) (_ gofhir.Patient, err error) {
	var patient gofhir.Patient
	dsnap, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if err != nil {
		return gofhir.Patient{}, err
	}
	dsnap.DataTo(&patient)
	return patient, nil
}

func (s *Service) Update(ctx context.Context) (err error) {
	return nil
}

func (s *Service) Delete(ctx context.Context, id string) (err error) {
	_, err = s.client.Collection(s.collection).Doc(id).Delete(ctx)
	return err
}
