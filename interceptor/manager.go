package interceptor

import (
	"reflect"

	"github.com/gabriel-ross/gofhir/hook"
)

// TODO: methods for registering new interceptors
//   - probably needs to be integrated into endpoint service
//
// TODO: how should this call/be called?
//
// Manager is a struct for managing and calling interceptors.
type Manager struct {
	Interceptors               []*interface{}
	ServerStartupInterceptors  []*hook.ServerStartupInterceptor
	ServerShutdownInterceptors []*hook.ServerShutdownInterceptor
	RequestInterceptors        []*hook.RequestInterceptor
	ResponseInterceptors       []*hook.ResponseInterceptor
	DatabaseInterceptors       []*hook.DatabaseInterceptor
	ErrorInterceptors          []*hook.ErrorInterceptor
}

func (m *Manager) RegisterInterceptor(i interface{}) {
	if interceptor, ok := i.(hook.RequestInterceptor); ok {
		m.RequestInterceptors = append(m.RequestInterceptors, &interceptor)
	}
	if interceptor, ok := i.(hook.ResponseInterceptor); ok {
		m.ResponseInterceptors = append(m.ResponseInterceptors, &interceptor)
	}
	if interceptor, ok := i.(hook.DatabaseInterceptor); ok {
		m.DatabaseInterceptors = append(m.DatabaseInterceptors, &interceptor)
	}
}

func (m *Manager) OnServerStartup(e *hook.ServerEvent)            {}
func (m *Manager) OnServerShutdown(e *hook.ServerEvent)           {}
func (m *Manager) OnRequestReceived(e *hook.RequestEvent)         {}
func (m *Manager) OnServerResponse(e *hook.ResponseEvent)         {}
func (m *Manager) BeforeDatabaseQuery(e *hook.DatabaseQueryEvent) {}
func (m *Manager) AfterDatabaseQuery(e *hook.DatabaseQueryEvent)  {}
func (m *Manager) OnError(e *hook.ErrorEvent)                     {}

var _ hook.Interceptor = &Manager{}

var fooIface interface{} = &Manager{}
var fooVal = reflect.ValueOf(fooIface).Elem()
var _, ok = fooIface.(hook.Interceptor) // TODO: Do something like this for registration, Gabe

// func (h *Hooker) Create(ctx context.Context) {
// 	event := DatabaseQueryEvent{}
// 	for _, interceptor := range h.BeforeDatabaseQueryInterceptors {
// 		interceptor.BeforeDatabaseQuery(&event)
// 		if event.abort == true {
// 			return
// 		}
// 	}

// 	// TODO: implement storage here directly

// 	for _, interceptor := range h.AfterDatabaseQueryInterceptors {
// 		interceptor.AfterDatabaseQuery(&event)
// 		if event.abort == true {
// 			return
// 		}
// 	}
// }
