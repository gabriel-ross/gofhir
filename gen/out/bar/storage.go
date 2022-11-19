package bar

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

func NewBarClient(c *firestore.Client) *StorageService {
	return &StorageService{
		client:     c,
		collection: "bar",
	}
}

func (s *StorageService) Create(ctx context.Context, bar gofhir.Bar) (_ gofhir.Bar, err error) {
	ref := s.client.Collection(s.collection).NewDoc()
	bar.ID = ref.ID
	_, err = s.client.Collection(s.collection).Doc(ref.ID).Set(ctx, bar)
	return bar, err
}

func (s *StorageService) List(ctx context.Context) (_ []gofhir.Bar, err error) {
	resp := []gofhir.Bar{}
	iter := s.client.Collection(s.collection).Documents(ctx)
	for {
		dsnap, done := iter.Next()
		if done == iterator.Done {
			break
		}
		var bar gofhir.Bar
		dsnap.DataTo(&bar)
		resp = append(resp, bar)
	}
	return resp, nil
}

func (s *StorageService) Read(ctx context.Context, id string) (_ gofhir.Bar, err error) {
	var bar gofhir.Bar
	dsnap, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if err != nil {
		return gofhir.Bar{}, err
	}
	dsnap.DataTo(&bar)
	return bar, nil
}

func (s *StorageService) Update(ctx context.Context) (err error) {
	return nil
}

func (s *StorageService) Delete(ctx context.Context, id string) (err error) {
	_, err = s.client.Collection(s.collection).Doc(id).Delete(ctx)
	return err
}
