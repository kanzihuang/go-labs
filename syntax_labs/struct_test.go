package syntax_labs

import (
	"testing"
)

// 空结构体 struct{}{} 不分配新的内存，复用已创建的空结构体
func TestStructEmpty(t *testing.T) {
	a := &struct{}{}
	b := &struct{}{}
	t.Logf("a: %p, b: %p\n", a, b)
	if a != b {
		t.Errorf("want: a == b, but got %p, %p\n", a, b)
	}
}
