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
	i.RequestTime[e.RequestID] = time.Now()
}

func (i *BenchmarkInterceptor) OnServerResponse(e *hook.ResponseEvent) {
	i.ResponseTime[e.RequestID] = time.Now()
	i.Avg()
}

func (i *BenchmarkInterceptor) BeforeDatabaseQuery(e *hook.DatabaseQueryEvent) {
	i.PreDbQueryTime[e.RequestID] = time.Now()
}
func (i *BenchmarkInterceptor) AfterDatabaseQuery(e *hook.DatabaseQueryEvent) {
	i.PostDbQueryTime[e.RequestID] = time.Now()
}

func (i *BenchmarkInterceptor) Avg() {
	var a float64
	for key, val := range i.ResponseTime {
		a += float64(val.Sub(i.RequestTime[key]).Milliseconds())
	}
	a = a / float64(len(i.ResponseTime))
	fmt.Printf("Request-response benchmarks:\nRequest volume: %d\nAverage latency: %vms\n", len(i.ResponseTime), a)

	var b float64
	for key, val := range i.PostDbQueryTime {
		b += float64(val.Sub(i.PreDbQueryTime[key]).Milliseconds())
	}
	b = b / float64(len(i.PostDbQueryTime))
	fmt.Printf("Database query benchmarks:\nRequest volume: %d\nAverage latency: %vms\n", len(i.PostDbQueryTime), b)
	fmt.Printf("Average server throughput: %v requests/s\n", (1.00/a)*1000.00)

}

var _ hook.RequestInterceptor = &BenchmarkInterceptor{}
var _ hook.ResponseInterceptor = &BenchmarkInterceptor{}
var _ hook.DatabaseInterceptor = &BenchmarkInterceptor{}
