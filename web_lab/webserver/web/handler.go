package web

import (
	"net/http"
)

type Handler struct {
	router Router
}

func NewHandler() *Handler {
	return &Handler{
		router: NewRouterBasedOnTree(),
	}
}

func (h *Handler) Route(method string, pattern string, handlerFunc HandlerFunc) {
	h.router.Route(method, pattern, handlerFunc)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handlerFunc := h.router.FindHandlerFunc(r.Method, r.URL.Path); handlerFunc != nil {
		c := NewContext(w, r)
		handlerFunc(c)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not any router matched"))
	}
}
