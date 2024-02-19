package simplerpc

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
)

type Server struct {
	proxy   Proxy
	startup chan struct{}
}

func NewServer(proxy Proxy) *Server {
	return &Server{
		proxy:   proxy,
		startup: make(chan struct{}),
	}
}

func (m *Server) WaitForStartup(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-m.startup:
		return nil
	}
}

func (m *Server) Start(network, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	m.startup <- struct{}{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("simplerpc: %v", err)
			continue
		}
		go m.handleConn(conn)
	}
}

func (m *Server) handleConn(conn net.Conn) {
	var err error
	req := &ProxyReq{}
	resp := &ProxyResp{}
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	for {
		err = decoder.Decode(req)
		if err != nil {
			slog.Error("simplerpc: %v", err)
			break
		}
		resp = m.proxy.Call(req)
		err = encoder.Encode(resp)
		if err != nil {
			slog.Error("simplerpc: %v", err)
			break
		}
	}
}
