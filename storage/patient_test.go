package storage

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/gofhir"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Patient = gofhir.Patient{
	Name: "test_name",
}

func TestCRUD(t *testing.T) {
	ctx := context.TODO()
	var err error
	c, err := firestore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	assert.Nil(t, err)
	defer c.Close()
	client := NewPatientClient(c)

	_, err = client.Create(ctx, Patient)
	assert.Nil(t, err)

	_, err = client.List(ctx)
	assert.Nil(t, err)

	err = client.Delete(ctx, Patient.ID)
	assert.Nil(t, err)

	_, err = client.Read(ctx, Patient.ID)
	assert.Equal(t, status.Code(err), codes.NotFound)
}
