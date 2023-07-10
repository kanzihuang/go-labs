package buildin_labs

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"testing"
)

func myFunc1() int {
	return 0
}

func myFunc2() int {
	return 0
}

func nameOf(fn func() int) string {
	if fn == nil {
		return "nil"
	}
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func TestFuncEqual(t *testing.T) {
	assert.NotEqual(t, nameOf(myFunc1), nameOf(nil))
	assert.Equal(t, nameOf(myFunc1), nameOf(myFunc1))
	assert.NotEqual(t, nameOf(myFunc1), nameOf(myFunc2))
	assert.Equal(t, runFuncName(), "go-labs/builtin_labs.TestFuncEqual")
}

// 获取正在运行的函数名
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
