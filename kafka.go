package gofhir

import (
	"context"
	"crypto/rand"
	"encoding/base32"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter(brokerAddress string) *KafkaWriter {
	return &KafkaWriter{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokerAddress),
			AllowAutoTopicCreation: true,
			Balancer:               &kafka.LeastBytes{},
		},
	}
}

func (w *KafkaWriter) WriteMessage(ctx context.Context, topic string, body []byte) (err error) {
	err = w.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(genToken(10)),
		Value: body,
	})
	if err != nil {
		return err
	}
	return nil
}

func (w *KafkaWriter) Close() (err error) {
	if err = w.writer.Close(); err != nil {
		return err
	}
	return nil
}

func genToken(length int) string {
	t := make([]byte, length)
	rand.Read(t)
	return base32.StdEncoding.EncodeToString(t)
}
