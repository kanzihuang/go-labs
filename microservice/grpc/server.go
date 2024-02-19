package grpc

import (
	"context"
	"google.golang.org/grpc"
	"net"
)

func NewServer(grpcSvr *grpc.Server) *Server {
	return &Server{
		startup: make(chan struct{}),
		grpcSvr: grpcSvr,
	}
}

type Server struct {
	startup  chan struct{}
	listener net.Listener
	grpcSvr  *grpc.Server
}

func (s *Server) Start(network, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	s.listener = listener
	s.startup <- struct{}{}
	return s.grpcSvr.Serve(listener)
}

func (s *Server) WaitForStartup(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.startup:
		return nil
	}
}

func (s *Server) Address() string {
	return s.listener.Addr().String()
}
