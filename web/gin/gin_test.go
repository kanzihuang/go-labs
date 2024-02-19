package gin

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGin_Ping(t *testing.T) {
	g := gin.Default()
	g.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	listener, err := net.Listen("tcp", ":")
	require.NoError(t, err)
	address := listener.Addr().String()
	go func() {
		err = http.Serve(listener, g.Handler())
		require.NoError(t, err)
	}()

	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/ping", address), nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	b := bytes.Buffer{}
	_, err = b.ReadFrom(resp.Body)
	require.NoError(t, err)
	err = resp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, "pong", b.String())
}

func TestGin_ServeHTTP(t *testing.T) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	require.NoError(t, err)
	r.ServeHTTP(rw, req)
	resp := rw.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, resp.Body.Close())
	require.NoError(t, err)
	require.Equal(t, []byte("pong"), body)
}

func TestGin_NewTLSServer(t *testing.T) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	ts := httptest.NewTLSServer(r.Handler())
	defer ts.Close()
	client := ts.Client()
	resp, err := client.Get(ts.URL + "/ping")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, resp.Body.Close())
	require.NoError(t, err)
	require.Equal(t, []byte("pong"), body)
}
