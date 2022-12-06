package interceptor

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gabriel-ross/gofhir"
	"github.com/gabriel-ross/gofhir/hook"
)

type StreamingInterceptor struct {
	KafkaWriter gofhir.KafkaWriter
}

func (i *StreamingInterceptor) BeforeDatabaseQuery(e *hook.DatabaseQueryEvent) {}

func (i *StreamingInterceptor) AfterDatabaseQuery(e *hook.DatabaseQueryEvent) {
	if e.Successful {
		bytes, err := json.Marshal(e)
		if err != nil {
			log.Println("error marshaling database event")
			return
		}
		i.KafkaWriter.WriteMessage(context.TODO(), "dbevent", bytes)
	}
}
