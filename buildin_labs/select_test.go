package buildin_labs

import (
	"testing"
	"time"
)

func TestSelectNil(t *testing.T) {
	var c chan int
	select {
	case <-c:
		t.Error("本处不应执行")
	default:
	}
	// all goroutines are asleep - deadlock!
	select {
	case <-c:
		t.Error("本处不应执行")
	}
	t.Error("本处不应执行")
}

func TestSelectTimeout(t *testing.T) {
	timeout := time.Second * 1
	select {
	case <-time.After(timeout):
		t.Log("Timeout")
		return
	}
	t.Error("本处不应执行")
}
