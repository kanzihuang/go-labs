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
	proxy.Register("Health", NewHealth())
	health := NewHealthClient(proxy)
	const length = 16
	payload := make([]byte, 0, length)
	for i := byte(0); i < length; i++ {
		payload = append(payload, i)
	}
	respPayload, err := health.Ping(payload)
	require.NoError(t, err)
	require.Equal(t, payload, respPayload)
}

func TestClientProxy(t *testing.T) {
	const network = "unix"
	address := filepath.Join(os.TempDir(), "simplerpc_test.sock")
	_ = os.Remove(address)
	proxy := NewServerProxy()
	err := proxy.Register("Health", NewHealth())
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
	respPayload, err := health.Ping(payload)
	require.NoError(t, err)
	require.Equal(t, payload, respPayload)
}
