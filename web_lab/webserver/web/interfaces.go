package web

import (
	"net/http"
)

// Server 是http server 的顶级抽象
type Server interface {
	Routable
	Start(address string) error
}

type Routable interface {
	Route(method string, pattern string, handlerFunc func(c *Context))
}

type Handler interface {
	http.Handler
	Routable
}
