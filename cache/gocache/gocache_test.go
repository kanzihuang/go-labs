package gocache

import (
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGoCache(t *testing.T) {
	c := cache.New(time.Minute, time.Hour)
	c.Set("foo", "bar", cache.DefaultExpiration)
	foo, found := c.Get("foo")
	assert.True(t, found)
	assert.Equal(t, "bar", foo)
}
