package web

import (
	"net/http"
)

type sdkHttpServer struct {
	Name   string
	router Router
	root   Filter
}

func (s *sdkHttpServer) Route(method, pattern string, handlerFunc HandlerFunc) {
	s.router.Route(method, pattern, handlerFunc)
}

func (s *sdkHttpServer) Start(address string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r)
		s.root(c)
	})
	return http.ListenAndServe(address, nil)
}

func (s *sdkHttpServer) ServeHTTP(c *Context) {
	if handlerFunc := s.router.FindHandlerFunc(c.R.Method, c.R.URL.Path); handlerFunc != nil {
		handlerFunc(c)
	} else {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("Not any router matched"))
	}
}

func NewServer(name string, builders ...FilterBuilder) Server {
	server := &sdkHttpServer{
		Name:   name,
		router: NewRouterBasedOnTree(),
	}
	var root Filter = server.ServeHTTP
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}
	server.root = root
	return server
}
