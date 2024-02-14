package silenceper

import (
	"github.com/silenceper/pool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

import (
	_ "github.com/silenceper/pool"
)

func TestPool(t *testing.T) {
	const (
		initialCap = 10
		maxIdle    = 30
		maxCap     = 50
	)
	config := &pool.Config{
		InitialCap: initialCap,
		MaxIdle:    maxIdle,
		MaxCap:     maxCap,
		Factory: func() (interface{}, error) {
			return net.Dial("tcp", ":3306")
		},
		Close: func(conn interface{}) error {
			return conn.(net.Conn).Close()
		},
	}
	p, err := pool.NewChannelPool(config)
	require.NoError(t, err)
	assert.Equal(t, initialCap, p.Len())

	conn, err := p.Get()
	require.NoError(t, err)
	assert.Equal(t, initialCap-1, p.Len())

	err = p.Put(conn)
	assert.NoError(t, err)
	assert.Equal(t, initialCap, p.Len())

	p.Release()
	assert.Equal(t, 0, p.Len())
}
