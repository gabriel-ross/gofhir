package interceptor

import (
	"fmt"

	"github.com/gabriel-ross/gofhir/hook"
)

type DemoInterceptor struct{}

func (i *DemoInterceptor) OnRequestReceived(e *hook.RequestEvent) {
	fmt.Println("hello world from the demo interceptor")
}

func (i *DemoInterceptor) OnServerResponse(e *hook.ResponseEvent) {
	fmt.Println("hello demo presentation from the response interceptor")
}

demo := &interceptor.DemoInterceptor{}