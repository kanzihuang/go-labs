package health

import (
	"context"
	"go-labs/microservice/grpc/health/gen"
)

type Service struct {
	gen.UnimplementedHealthServiceServer
}

func (s *Service) Ping(_ context.Context, req *gen.PingReq) (*gen.PingResp, error) {
	resp := &gen.PingResp{
		Payload: req.Payload,
	}
	return resp, nil
}
