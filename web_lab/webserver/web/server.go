package web

import (
	"net/http"
)

type sdkHttpServer struct {
	Name    string
	handler *Handler
}

func (s *sdkHttpServer) Route(method, pattern string, handlerFunc handlerFunc) {
	s.handler.Route(method, pattern, handlerFunc)
}

func (s *sdkHttpServer) Start(address string) error {
	return http.ListenAndServe(address, s.handler)
}

func NewServer() Server {
	return &sdkHttpServer{
		Name:    "sdkHttpServer",
		handler: NewHandler(),
	}
}
