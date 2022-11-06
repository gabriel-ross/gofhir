package gofhir

import (
	"context"

	"cloud.google.com/go/logging"
)

func NewCloudLoggerClient(ctx context.Context, projectID string) (_ *logging.Client, err error) {
	return logging.NewClient(ctx, projectID)
}
