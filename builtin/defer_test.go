package builtin

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefer(t *testing.T) {
	const (
		deferName  = "deferName"
		returnName = "returnName"
	)
	getName := func() (name string) {
		defer func() {
			name = deferName
		}()
		return returnName
	}
	name := getName()
	require.Equal(t, deferName, name)
}
