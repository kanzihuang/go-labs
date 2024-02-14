package builtin

import (
	"github.com/stretchr/testify/assert"
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

type StructTest struct {
	num   int
	names []string
}

func copyStruct(s StructTest) StructTest {
	return s
}

func TestCopyStruct(t *testing.T) {
	a := StructTest{
		num:   1,
		names: []string{"Apple"},
	}
	b := copyStruct(a)
	assert.Equal(t, 1, b.num)
	assert.Equal(t, 1, len(b.names))

	a.num = -1
	assert.Equal(t, 1, b.num)

	a.names[0] = "Tomato"
	assert.Equal(t, "Tomato", b.names[0])

	a.names = []string{"potato"}
	assert.Equal(t, "Tomato", b.names[0])
}
