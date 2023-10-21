package builtin

import (
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestRangeInt(t *testing.T) {
	a := []int{1, 2, 3}
	want := len(a)
	i := 0
	for range a {
		log.Println("index:", i)
		a = append(a, i+4)
		i++
	}
	if i != want {
		t.Errorf("index after range got: %d, want:%d ", i, want)
	}
}

func TestRangeChan(t *testing.T) {
	tm := time.After(time.Millisecond * 100)
	for range tm {
		t.Log("Timeout")
		break
	}

}

func TestCopy(t *testing.T) {
	hello := []byte("hello world")
	copy(hello, hello[3:3])
	require.Equal(t, "hello world", string(hello))
}
