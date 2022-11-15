package patient

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/gofhir"
	"google.golang.org/api/iterator"
)

type firestoreService struct {
	client     *firestore.Client
	collection string
}

func NewFirestoreService(c *firestore.Client) *firestoreService {
	return &firestoreService{
		client:     c,
		collection: "patient",
	}
}

func (c *firestoreService) Create(ctx context.Context, patient gofhir.Patient) (_ gofhir.Patient, err error) {
	ref := c.client.Collection(c.collection).NewDoc()
	patient.ID = ref.ID
	_, err = c.client.Collection(c.collection).Doc(ref.ID).Set(ctx, patient)
	return patient, err
}

func (c *firestoreService) List(ctx context.Context, options ...gofhir.ListOption) (_ []gofhir.Patient, err error) {
	resp := []gofhir.Patient{}
	query := gofhir.BuildListQueryFromListOptions(c.client.Collection(c.collection).Query, options...)
	iter := query.Documents(ctx)
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

func (c *firestoreService) Read(ctx context.Context, id string) (_ gofhir.Patient, err error) {
	var patient gofhir.Patient
	dsnap, err := c.client.Collection(c.collection).Doc(id).Get(ctx)
	if err != nil {
		return gofhir.Patient{}, err
	}
	dsnap.DataTo(&patient)
	return patient, nil
}

func (c *firestoreService) Update(ctx context.Context, id string, patient gofhir.Patient) (_ gofhir.Patient, err error) {
	_, err = c.client.Collection(c.collection).Doc(id).Set(ctx, patient)
	return patient, err
}

func (c *firestoreService) Delete(ctx context.Context, id string) (err error) {
	_, err = c.client.Collection(c.collection).Doc(id).Delete(ctx)
	return err
}
