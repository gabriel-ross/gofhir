package interceptor

import (
	"fmt"
	"time"

	"github.com/gabriel-ross/gofhir/hook"
)

type BenchmarkInterceptor struct {
	RequestTime     map[string]time.Time
	PreDbQueryTime  map[string]time.Time
	PostDbQueryTime map[string]time.Time
	ResponseTime    map[string]time.Time
}

func (i *BenchmarkInterceptor) OnRequestReceived(e *hook.RequestEvent) {
	i.RequestTime[e.RequestID] = e.Timestamp
}

func (i *BenchmarkInterceptor) OnServerResponse(e *hook.ResponseEvent) {
	i.ResponseTime[e.RequestID] = e.Timestamp

	if len(i.ResponseTime)%10 == 0 {
		i.Avg()
	}
}

func (i *BenchmarkInterceptor) BeforeDatabaseQuery(e *hook.DatabaseQueryEvent) {
	i.PreDbQueryTime[e.RequestID] = e.Timestamp
}
func (i *BenchmarkInterceptor) AfterDatabaseQuery(e *hook.DatabaseQueryEvent) {
	i.PostDbQueryTime[e.RequestID] = e.Timestamp
}

func (i *BenchmarkInterceptor) Avg() {
	var a int64
	for key, val := range i.ResponseTime {
		a += val.Sub(i.RequestTime[key]).Milliseconds()
	}
	a = a / int64(len(i.ResponseTime))
	fmt.Printf("Request-response benchmarks:\nRequest volume: %d\nAverage latency: %d\n\n", len(i.ResponseTime), a)

	var b int64
	for key, val := range i.PostDbQueryTime {
		b += val.Sub(i.PreDbQueryTime[key]).Milliseconds()
	}
	b = b / int64(len(i.PostDbQueryTime))
	fmt.Printf("Database query benchmarks:\nRequest volume: %d\nAverage latency: %d\n\n", len(i.PostDbQueryTime), b)

	fmt.Printf("Request breakdown:\nDatabase operations: %d%\nBusiness logic & interceptors: %s%\n\n", b/a, (a-b)/a)
	fmt.Printf("Server throughput over the last 10 requests: %s requests/s\n", (1/a)*1000)

}

var _ hook.RequestInterceptor = &BenchmarkInterceptor{}
var _ hook.ResponseInterceptor = &BenchmarkInterceptor{}
var _ hook.DatabaseInterceptor = &BenchmarkInterceptor{}
