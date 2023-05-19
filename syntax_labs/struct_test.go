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

type Node struct {
	value int
	left  *Node
	right *Node
}

func setNext(next **Node, node *Node) {
	*next = node
}

func NewNode(value int) *Node {
	return &Node{value: value}
}

func TestPoint(t *testing.T) {
	root := NewNode(5)
	root.left = NewNode(2)
	root.right = NewNode(10)
	setNext(&root.right, nil)
}
