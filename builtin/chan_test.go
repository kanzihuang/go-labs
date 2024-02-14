package builtin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func copyChan(c chan struct{}) chan struct{} {
	return c
}

func TestChanLen(t *testing.T) {
	c1 := make(chan struct{}, 1)
	c1 <- struct{}{}
	assert.Equal(t, 1, len(c1))

	c2 := copyChan(c1)
	assert.Equal(t, 1, len(c2))

	<-c1
	assert.Equal(t, 0, len(c1))
	assert.Equal(t, 0, len(c2))
}
