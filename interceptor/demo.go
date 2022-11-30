package interceptor

import (
	"fmt"
	"time"

	"github.com/gabriel-ross/gofhir/hook"
)

type DemoInterceptor struct {
	RequestTime     map[string]time.Time
	RequestDuration map[string]int64
	DbQueryTime     map[string]time.Time
	DbQueryDuration map[string]int64
}

func (i *DemoInterceptor) OnRequestReceived(e *hook.RequestEvent) {
	i.RequestTime[e.RequestID] = e.Timestamp
}

func (i *DemoInterceptor) OnServerResponse(e *hook.ResponseEvent) {
	if reqTimeStamp, ok := i.RequestTime[e.RequestID]; ok {
		i.RequestDuration[e.RequestID] = time.Since(reqTimeStamp).Milliseconds()
		fmt.Printf("Response event\nid: %s\ntime since request: %v\n\n", e.RequestID, time.Since(reqTimeStamp))
	}

	if len(i.RequestDuration)%10 == 0 {
		i.Avg()
	}
}

func (i *DemoInterceptor) BeforeDatabaseQuery(e *hook.DatabaseQueryEvent) {
	i.DbQueryTime[e.RequestID] = e.Timestamp
}
func (i *DemoInterceptor) AfterDatabaseQuery(e *hook.DatabaseQueryEvent) {
	if reqTimeStamp, ok := i.DbQueryTime[e.RequestID]; ok {
		i.DbQueryDuration[e.RequestID] = time.Since(reqTimeStamp).Milliseconds()
		fmt.Printf("After db query event\nid: %s\ntime since query: %v\n\n", e.RequestID, time.Since(reqTimeStamp))
	}
}

func (i *DemoInterceptor) Avg() {
	var reqAvg int64
	for _, val := range i.RequestDuration {
		reqAvg += val
	}
	reqAvg = reqAvg / int64(len(i.RequestDuration))

	var dbQueryAvg int64
	for _, val := range i.DbQueryDuration {
		dbQueryAvg += val
	}
	dbQueryAvg = dbQueryAvg / int64(len(i.DbQueryDuration))

	fmt.Printf("Request-response benchmarks:\nRequest volume: %d\nAverage latency: %d\n\n", len(i.RequestDuration), reqAvg)
	fmt.Printf("Database query benchmarks:\nRequest volume: %d\nAverage latency: %d\n\n", len(i.DbQueryDuration), dbQueryAvg)
}

var _ hook.RequestInterceptor = &DemoInterceptor{}
var _ hook.ResponseInterceptor = &DemoInterceptor{}
var _ hook.DatabaseInterceptor = &DemoInterceptor{}
