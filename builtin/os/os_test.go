package os

import (
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestPipe(t *testing.T) {
	r, w, err := os.Pipe()
	require.NoError(t, err)
	//out := strings.Builder{}
	//go func() {
	_, err = w.Write([]byte("hello"))
	require.NoError(t, err)
	err = w.Close()
	require.NoError(t, err)
	//}()
	out, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), out)
}
