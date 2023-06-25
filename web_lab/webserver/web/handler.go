package web

import (
	"net/http"
)

type Handler struct {
	router Router
}

func NewHandler() *Handler {
	return &Handler{
		router: NewRouterBaseOnMap(),
	}
}

func (h *Handler) Route(method string, pattern string, handlerFunc handlerFunc) {
	h.router.Route(method, pattern, handlerFunc)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{W: w, R: r}
	if h.router.handle(r.Method, r.URL.Path, c) != true {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not any router matched"))
	}
}
