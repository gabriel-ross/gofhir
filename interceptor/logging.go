package interceptor

import (
	"fmt"
	"log"

	"github.com/gabriel-ross/gofhir"
	"github.com/gabriel-ross/gofhir/hook"
)

type LoggingInterceptor struct {
	Logger *gofhir.CloudLogger
}

func (i *LoggingInterceptor) OnRequestReceived(e *hook.RequestEvent) {
	entry := fmt.Sprintf("RequestID: %s\nReceived at: %v\nRequestor: %s", e.RequestID, e.Timestamp, e.Request.RemoteAddr)
	i.Logger.Log(entry)
	log.Println(entry)
}
