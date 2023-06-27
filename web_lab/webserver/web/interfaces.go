package web

type HandlerFunc func(c *Context)

// Server 是http server 的顶级抽象
type Server interface {
	Routable
	Start(address string) error
}

type Routable interface {
	Route(method string, pattern string, handlerFunc HandlerFunc)
}

type Router interface {
	Routable
	//handle(method string, path string, context *Context) bool
	FindHandlerFunc(method string, path string) HandlerFunc
}
