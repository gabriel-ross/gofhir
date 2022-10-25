package gofhir

import (
	"context"

	"cloud.google.com/go/firestore"
)

func NewFirestoreClient(ctx context.Context, projectID string) (_ *firestore.Client, err error) {
	return firestore.NewClient(ctx, projectID)
}
