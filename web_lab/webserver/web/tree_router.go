package web

import "strings"

type Node interface {
	FindNode(method, path string) Node
	AddNode(method string, pattern string, handlerFunc handlerFunc)
	Handle(context *Context) bool
	matchNames(method string, names []string) Node
	isMatched(method string, name string) bool
	addNodeByNames(method string, names []string, handlerFunc handlerFunc)
	isRoot() bool
}

type RouterBasedOnTree struct {
	root Node
}

func NewRouterBasedOnTree() Router {
	return &RouterBasedOnTree{
		root: NewExactNode("", "", nil),
	}
}

type BaseNode struct {
	method  string
	name    string
	handler handlerFunc
}

type ExactNode struct {
	BaseNode
	children []Node
}

func (n *ExactNode) isRoot() bool {
	return n.method == ""
}

func NewExactNode(method string, name string, handlerFunc handlerFunc) Node {
	return &ExactNode{
		BaseNode: BaseNode{
			method:  method,
			name:    name,
			handler: handlerFunc,
		},
	}
}

func (n *ExactNode) addNodeByNames(method string, names []string, handlerFunc handlerFunc) {

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
		node.addNodeByNames(method, names[1:], handlerFunc)
	}
}

func (n *ExactNode) isMatched(method string, name string) bool {
	return n.method == method && n.name == name
}

func (n *ExactNode) matchNames(method string, names []string) Node {
	var isMatched bool
	if n.isRoot() {
		isMatched = true
	} else {
		isMatched = n.isMatched(method, names[0])
		names = names[1:]
	}
	if len(names) == 0 {
		if isMatched {
			return n
		} else {
			return nil
		}
	}
	if isMatched {
		for _, child := range n.children {
			if matched := child.matchNames(method, names); matched != nil {
				return matched
			}
		}
	}
	return nil
}

func (n *ExactNode) FindNode(method, path string) Node {
	path = strings.TrimRight(path, "/")
	names := strings.Split(path, "/")
	return n.matchNames(method, names)
}

func (n *ExactNode) AddNode(method string, pattern string, handlerFunc handlerFunc) {
	pattern = strings.TrimRight(pattern, "/")
	names := strings.Split(pattern, "/")
	n.addNodeByNames(method, names, handlerFunc)
}

func (n *ExactNode) Handle(context *Context) bool {
	if n.handler != nil {
		n.handler(context)
		return true
	} else {
		return false
	}
}

func (h *RouterBasedOnTree) Route(method string, pattern string, handlerFunc handlerFunc) {
	h.root.AddNode(method, pattern, handlerFunc)
}

func (h *RouterBasedOnTree) handle(method string, path string, context *Context) bool {
	if node := h.root.FindNode(method, path); node != nil {
		return node.Handle(context)
	}
	return false
}
