package gofhir

import (
	"context"

	"cloud.google.com/go/logging"
)

type CloudLogger struct {
	client *logging.Client
	logger *logging.Logger
}

func NewCloudLogger(ctx context.Context, projectID string, logName string) (_ *CloudLogger, err error) {
	clClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &CloudLogger{
		client: clClient,
		logger: clClient.Logger(logName),
	}, nil
}

func (cl *CloudLogger) Log(body string) {
	cl.logger.Log(logging.Entry{
		Payload: body,
	})
}

func (cl *CloudLogger) Close() (err error) {
	if err = cl.client.Close(); err != nil {
		return err
	}
	return nil
}
