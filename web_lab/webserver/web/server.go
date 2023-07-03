package web

import (
	"net/http"
	"sync"
)

type sdkHttpServer struct {
	Router
	Name        string
	root        Filter
	contextPool sync.Pool
}

func (s *sdkHttpServer) Start(address string) error {
	return http.ListenAndServe(address, s)
}

func (s *sdkHttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := s.contextPool.Get().(*Context)
	defer s.contextPool.Put(c)
	c.Reset(w, r)
	s.root(c)
}

func (s *sdkHttpServer) ServeHTTPWithContext(c *Context) {
	if handlerFunc := s.FindHandlerFunc(c.R.Method, c.R.URL.Path); handlerFunc != nil {
		handlerFunc(c)
	} else {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("Not any router matched"))
	}
}

func NewServer(name string, builders ...FilterBuilder) Server {
	server := &sdkHttpServer{
		Router: NewRouterBasedOnTree(),
		Name:   name,
		contextPool: sync.Pool{
			New: NewContext,
		},
	}
	var root Filter = server.ServeHTTPWithContext
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}
	server.root = root
	return server
}
