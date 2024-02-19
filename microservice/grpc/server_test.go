package grpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-labs/microservice/grpc/health"
	"go-labs/microservice/grpc/health/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	const svrNum = 3
	addresses := make([]string, 0, svrNum)
	servers := make([]*Server, 0, svrNum)
	for i := 0; i < svrNum; i++ {
		grpcSvr := grpc.NewServer()
		gen.RegisterHealthServiceServer(grpcSvr, &health.Service{})
		svr := NewServer(grpcSvr)
		go func() {
			_ = svr.Start("tcp", ":")
		}()
		err := svr.WaitForStartup(ctx)
		require.NoError(t, err)
		addresses = append(addresses, svr.Address())
		servers = append(servers, svr)
	}

	conn, err := grpc.Dial("registrar://",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(NewResolverBuilder(addresses...)))
	require.NoError(t, err)
	defer func() {
		err = conn.Close()
		require.NoError(t, err)
	}()

	client := gen.NewHealthServiceClient(conn)
	req := &gen.PingReq{
		Payload: "abcdefghijklmn",
	}
	resp, err := client.Ping(ctx, req)
	require.NoError(t, err)
	require.Equal(t, req.Payload, resp.GetPayload())
}
