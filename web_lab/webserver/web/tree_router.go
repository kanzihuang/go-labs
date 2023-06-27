package web

import "strings"

type Node interface {
	findChild(method string, name string) Node
	findHandlerFunc(method string, names []string) HandlerFunc
	isMatched(method string, name string) bool
	addChild(method string, names []string, handlerFunc HandlerFunc)
	isRoot() bool
	getHandlerFunc() HandlerFunc
}

type RouterBasedOnTree struct {
	root Node
}

func (h *RouterBasedOnTree) FindHandlerFunc(method string, path string) HandlerFunc {
	path = strings.TrimRight(path, "/")
	names := strings.Split(path, "/")
	return h.root.findHandlerFunc(method, names)
}

func NewRouterBasedOnTree() Router {
	return &RouterBasedOnTree{
		root: NewExactNode("", "", nil),
	}
}

type BaseNode struct {
	method      string
	name        string
	handlerFunc HandlerFunc
}

type ExactNode struct {
	BaseNode
	children []Node
}

func (n *ExactNode) getHandlerFunc() HandlerFunc {
	return n.handlerFunc
}

func (n *ExactNode) isRoot() bool {
	return n.method == ""
}

func NewExactNode(method string, name string, handlerFunc HandlerFunc) Node {
	return &ExactNode{
		BaseNode: BaseNode{
			method:      method,
			name:        name,
			handlerFunc: handlerFunc,
		},
	}
}

func (n *ExactNode) addChild(method string, names []string, handlerFunc HandlerFunc) {

	var node Node = nil
	name := names[0]
	for _, child := range n.children {
		if child.isMatched(method, name) {
			node = child
			break
		}
	}
	if node == nil {
		if len(names) <= 1 {
			node = NewExactNode(method, name, handlerFunc)
		} else {
			node = NewExactNode(method, name, nil)
		}
		n.children = append(n.children, node)
	}
	if len(names) > 1 {
		node.addChild(method, names[1:], handlerFunc)
	}
}

func (n *ExactNode) isMatched(method string, name string) bool {
	return n.method == method && n.name == name
}

func (n *ExactNode) findChild(method string, name string) Node {
	for _, child := range n.children {
		if child.isMatched(method, name) {
			return child
		}
	}
	return nil
}

func (n *ExactNode) findHandlerFunc(method string, names []string) HandlerFunc {
	if matched := n.findChild(method, names[0]); matched != nil {
		if len(names) > 1 {
			if handlerFunc := matched.findHandlerFunc(method, names[1:]); handlerFunc != nil {
				return handlerFunc
			}
		}
		return matched.getHandlerFunc()
	}
	return n.getHandlerFunc()
}

func (h *RouterBasedOnTree) Route(method string, pattern string, handlerFunc HandlerFunc) {
	pattern = strings.TrimRight(pattern, "/")
	names := strings.Split(pattern, "/")
	h.root.addChild(method, names, handlerFunc)
}
