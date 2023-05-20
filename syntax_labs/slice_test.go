package syntax_labs

import (
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
	c := time.After(time.Second * 1)
	// all goroutines are asleep - deadlock!
	for range c {
		log.Println("timeout")
	}
	t.Error("此处不应执行")
}
