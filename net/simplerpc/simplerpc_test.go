package simplerpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestServerProxy(t *testing.T) {
	proxy := NewServerProxy()
	err := proxy.Register("Health", NewHealth())
	require.NoError(t, err)
	health := NewHealthClient(proxy)
	const length = 16
	payload := make([]byte, 0, length)
	for i := byte(0); i < length; i++ {
		payload = append(payload, i)
	}
	pingReq := &PingReq{
		Payload: payload,
	}
	pingResp, err := health.Ping(pingReq)
	require.NoError(t, err)
	require.Equal(t, payload, pingResp.Payload)
}

func testStartServer(t *testing.T, network, address string) {
	var err error
	_ = os.Remove(address)
	proxy := NewServerProxy()
	err = proxy.Register("Health", NewHealth())
	require.NoError(t, err)
	server := NewServer(proxy)
	go func() {
		err := server.Start(network, address)
		slog.Info("simplerpc: %v", err)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = server.WaitForStartup(ctx)
	require.NoError(t, err)
}

func TestClientProxy(t *testing.T) {
	var err error
	const network = "unix"
	address := filepath.Join(os.TempDir(), "simplerpc_test.sock")
	testStartServer(t, network, address)

	client, err := NewClient(network, address)
	require.NoError(t, err)
	defer func() {
		_ = client.Close()
	}()
	health := NewHealthClient(&ClientProxy{
		client: client,
	})
	const length = 16
	payload := make([]byte, 0, length)
	for i := byte(0); i < length; i++ {
		payload = append(payload, i)
	}
	pingReq := &PingReq{
		Payload: payload,
	}
	pingResp, err := health.Ping(pingReq)
	require.NoError(t, err)
	require.Equal(t, payload, pingResp.Payload)
}
