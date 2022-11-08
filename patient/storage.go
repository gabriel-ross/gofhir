package patient

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/gofhir"
	"google.golang.org/api/iterator"
)

type StorageService struct {
	client     *firestore.Client
	collection string
}

func NewPatientClient(c *firestore.Client) *StorageService {
	return &StorageService{
		client:     c,
		collection: "patient",
	}
}

func (s *StorageService) Create(ctx context.Context, patient gofhir.Patient) (_ gofhir.Patient, err error) {
	ref := s.client.Collection(s.collection).NewDoc()
	patient.ID = ref.ID
	_, err = s.client.Collection(s.collection).Doc(ref.ID).Set(ctx, patient)
	return patient, err
}

func (s *StorageService) List(ctx context.Context) (_ []gofhir.Patient, err error) {
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

func (s *StorageService) Read(ctx context.Context, id string) (_ gofhir.Patient, err error) {
	var patient gofhir.Patient
	dsnap, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if err != nil {
		return gofhir.Patient{}, err
	}
	dsnap.DataTo(&patient)
	return patient, nil
}

func (s *StorageService) Update(ctx context.Context) (err error) {
	return nil
}

func (s *StorageService) Delete(ctx context.Context, id string) (err error) {
	_, err = s.client.Collection(s.collection).Doc(id).Delete(ctx)
	return err
}
