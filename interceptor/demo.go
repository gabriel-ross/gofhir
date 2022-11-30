package interceptor

import (
	"fmt"
	"time"

	"github.com/gabriel-ross/gofhir/hook"
)

type DemoInterceptor struct {
	ReqLog     map[string]time.Time
	DbQueryLog map[string]time.Time
}

func (i *DemoInterceptor) OnRequestReceived(e *hook.RequestEvent) {
	i.ReqLog[e.RequestID] = e.Timestamp
	fmt.Printf("Request event\nid: %s\ntime: %v\n\n", e.RequestID, e.Timestamp)
}

func (i *DemoInterceptor) OnServerResponse(e *hook.ResponseEvent) {
	if reqTimeStamp, ok := i.ReqLog[e.RequestID]; ok {
		fmt.Printf("Response event\nid: %s\ntime since request: %v\n\n", e.RequestID, time.Since(reqTimeStamp))
	} else {
		fmt.Printf("Response event with no corresponding request id\nid: %s\n\n", e.RequestID)
	}
}

func (i *DemoInterceptor) BeforeDatabaseQuery(e *hook.DatabaseQueryEvent) {
	i.DbQueryLog[e.RequestID] = e.Timestamp
}
func (i *DemoInterceptor) AfterDatabaseQuery(e *hook.DatabaseQueryEvent) {
	if reqTimeStamp, ok := i.DbQueryLog[e.RequestID]; ok {
		fmt.Printf("After db query event\nid: %s\ntime since query: %v\n\n", e.RequestID, time.Since(reqTimeStamp))
	} else {
		fmt.Printf("After db query event with no corresponding before db query event id\nid: %s\n\n", e.RequestID)
	}
}

var _ hook.RequestInterceptor = &DemoInterceptor{}
var _ hook.ResponseInterceptor = &DemoInterceptor{}
var _ hook.DatabaseInterceptor = &DemoInterceptor{}
