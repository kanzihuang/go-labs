package web

import "strings"

type Node interface {
	findChild(method string, name string) Node
	findHandlerFunc(method string, names []string) HandlerFunc
	isMatched(method string, name string) bool
	addChild(method string, names []string, handlerFunc HandlerFunc)
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
	children map[string]Node
}

func (n *ExactNode) getHandlerFunc() HandlerFunc {
	return n.handlerFunc
}

func NewExactNode(method string, name string, handlerFunc HandlerFunc) Node {
	return &ExactNode{
		BaseNode: BaseNode{
			method:      method,
			name:        name,
			handlerFunc: handlerFunc,
		},
		children: make(map[string]Node),
	}
}

func (n *ExactNode) addChild(method string, names []string, handlerFunc HandlerFunc) {
	name := names[0]
	node := n.findChild(method, name)
	if node == nil {
		if len(names) <= 1 {
			node = NewExactNode(method, name, handlerFunc)
		} else {
			node = NewExactNode(method, name, nil)
		}
		n.children[n.key(method, name)] = node
	}
	if len(names) > 1 {
		node.addChild(method, names[1:], handlerFunc)
	}
}

func (n *ExactNode) isMatched(method string, name string) bool {
	return n.method == method && (n.name == "*" || n.name == name)
}

func (n *ExactNode) findChild(method string, name string) Node {
	if child := n.children[n.key(method, name)]; child != nil {
		return child
	} else if child := n.children[n.key(method, "*")]; child != nil {
		return child
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

func (n *ExactNode) key(method string, name string) string {
	return method + "#" + name
}

func (h *RouterBasedOnTree) Route(method string, pattern string, handlerFunc HandlerFunc) {
	pattern = strings.TrimRight(pattern, "/")
	names := strings.Split(pattern, "/")
	h.root.addChild(method, names, handlerFunc)
}
