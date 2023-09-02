package web

import "fmt"

type RouterBasedOnMap struct {
	handlers map[string]HandlerFunc
}

func (r *RouterBasedOnMap) FindHandlerFunc(method string, path string) HandlerFunc {
	key := r.key(method, path)
	return r.handlers[key]
}

func NewRouterBasedOnMap() Router {
	return &RouterBasedOnMap{
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *RouterBasedOnMap) handle(method string, path string, context *Context) bool {
	key := r.key(method, path)
	if handler, ok := r.handlers[key]; ok {
		handler(context)
		return true
	} else {
		return false
	}
}

func (r *RouterBasedOnMap) key(method string, path string) string {
	return fmt.Sprintf("%s#%s", method, path)
}

func (r *RouterBasedOnMap) Route(method string, pattern string, handlerFunc HandlerFunc) {
	key := r.key(method, pattern)
	r.handlers[key] = handlerFunc
}
