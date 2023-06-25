package web

import (
	"fmt"
	"net/http"
)

type HandlerBasedOnMap struct {
	handlers map[string]func(c *Context)
}

var _ Handler = (*HandlerBasedOnMap)(nil)

func NewHandlerBasedOnMap() Handler {
	return &HandlerBasedOnMap{
		handlers: map[string]func(c *Context){},
	}
}

func (h *HandlerBasedOnMap) Route(method string, pattern string, handlerFunc func(c *Context)) {
	key := h.key(method, pattern)
	h.handlers[key] = handlerFunc
}

func (h *HandlerBasedOnMap) key(method string, path string) string {
	return fmt.Sprintf("%s#%s", method, path)
}

func (h *HandlerBasedOnMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := h.key(r.Method, r.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		c := &Context{W: w, R: r}
		handler(c)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not any router matched"))
	}
}
