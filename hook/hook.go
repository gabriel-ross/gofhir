package hook

import "reflect"

var ServerStartupInterceptorType = reflect.TypeOf((*ServerStartupInterceptor)(nil)).Elem()
var ServerShutdownInterceptorType = reflect.TypeOf((*ServerShutdownInterceptor)(nil)).Elem()
var RequestInterceptorType = reflect.TypeOf((*RequestInterceptor)(nil)).Elem()
var ResponseInterceptorType = reflect.TypeOf((*ResponseInterceptor)(nil)).Elem()
var DatabaseInterceptorType = reflect.TypeOf((*DatabaseInterceptor)(nil)).Elem()
var ErrorInterceptorType = reflect.TypeOf((*ErrorInterceptor)(nil)).Elem()

type Interceptor interface {
	ServerStartupInterceptor
	ServerShutdownInterceptor
	RequestInterceptor
	ResponseInterceptor
	DatabaseInterceptor
	ErrorInterceptor
}

type ServerStartupInterceptor interface {
	OnServerStartup(*ServerEvent)
}

type ServerShutdownInterceptor interface {
	OnServerShutdown(*ServerEvent)
}

type RequestInterceptor interface {
	OnRequestReceived(*RequestEvent)
}

type ResponseInterceptor interface {
	OnServerResponse(*ResponseEvent)
}

type DatabaseInterceptor interface {
	BeforeDatabaseQuery(*DatabaseQueryEvent)
	AfterDatabaseQuery(*DatabaseQueryEvent)
}

type ErrorInterceptor interface {
	OnError(*ErrorEvent)
}
