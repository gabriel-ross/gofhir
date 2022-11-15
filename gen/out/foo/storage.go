package foo

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

func NewFooClient(c *firestore.Client) *StorageService {
	return &StorageService{
		client:     c,
		collection: "foo",
	}
}

func (s *StorageService) Create(ctx context.Context, foo gofhir.Foo) (_ gofhir.Foo, err error) {
	ref := s.client.Collection(s.collection).NewDoc()
	foo.ID = ref.ID
	_, err = s.client.Collection(s.collection).Doc(ref.ID).Set(ctx, foo)
	return foo, err
}

func (s *StorageService) List(ctx context.Context) (_ []gofhir.Foo, err error) {
	resp := []gofhir.Foo{}
	iter := s.client.Collection(s.collection).Documents(ctx)
	for {
		dsnap, done := iter.Next()
		if done == iterator.Done {
			break
		}
		var foo gofhir.Foo
		dsnap.DataTo(&foo)
		resp = append(resp, foo)
	}
	return resp, nil
}

func (s *StorageService) Read(ctx context.Context, id string) (_ gofhir.Foo, err error) {
	var foo gofhir.Foo
	dsnap, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if err != nil {
		return gofhir.Foo{}, err
	}
	dsnap.DataTo(&foo)
	return foo, nil
}

func (s *StorageService) Update(ctx context.Context) (err error) {
	return nil
}

func (s *StorageService) Delete(ctx context.Context, id string) (err error) {
	_, err = s.client.Collection(s.collection).Doc(id).Delete(ctx)
	return err
}
