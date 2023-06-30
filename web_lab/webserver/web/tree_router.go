package web

import (
	"log"
	"net/http"
	"strings"
)

type Node interface {
	findChild(name string) Node
	findHandlerFunc(names []string) HandlerFunc
	isMatched(name string) bool
	addChild(names []string, handlerFunc HandlerFunc)
	getHandlerFunc() HandlerFunc
	setHandlerFunc(handlerFunc HandlerFunc)
}

var supportMethods = [4]string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
}

type RouterBasedOnTree struct {
	forest map[string]Node
	root   Node
}

func (r *RouterBasedOnTree) Route(method string, pattern string, handlerFunc HandlerFunc) {
	pattern = strings.Trim(pattern, "/")
	names := strings.Split(pattern, "/")
	root, ok := r.forest[method]
	if ok {
		if len(names) == 1 && names[0] == "" {
			root.setHandlerFunc(handlerFunc)
		} else {
			root.addChild(names, handlerFunc)
		}
	} else {
		log.Panic("Not supported method", method)
	}
}

func (r *RouterBasedOnTree) FindHandlerFunc(method string, path string) HandlerFunc {
	path = strings.Trim(path, "/")
	names := strings.Split(path, "/")
	root, ok := r.forest[method]
	if ok {
		if len(names) == 1 && names[0] == "" {
			return root.getHandlerFunc()
		} else {
			if handlerFunc := root.findHandlerFunc(names); handlerFunc != nil {
				return handlerFunc
			} else {
				return root.getHandlerFunc()
			}
		}
	} else {
		return nil
	}
}

func NewRouterBasedOnTree() Router {
	forest := make(map[string]Node, len(supportMethods))
	for _, method := range supportMethods {
		forest[method] = NewExactNode("", nil)
	}
	return &RouterBasedOnTree{
		forest: forest,
	}
}

type BaseNode struct {
	name        string
	handlerFunc HandlerFunc
}

type ExactNode struct {
	BaseNode
	children map[string]Node
}

func (n *ExactNode) setHandlerFunc(handlerFunc HandlerFunc) {
	n.handlerFunc = handlerFunc
}

func (n *ExactNode) getHandlerFunc() HandlerFunc {
	return n.handlerFunc
}

func NewExactNode(name string, handlerFunc HandlerFunc) Node {
	return &ExactNode{
		BaseNode: BaseNode{
			name:        name,
			handlerFunc: handlerFunc,
		},
		children: make(map[string]Node),
	}
}

func (n *ExactNode) addChild(names []string, handlerFunc HandlerFunc) {
	name := names[0]
	node := n.findChild(name)
	if node == nil {
		if len(names) <= 1 {
			node = NewExactNode(name, handlerFunc)
		} else {
			node = NewExactNode(name, nil)
		}
		n.children[name] = node
	}
	if len(names) > 1 {
		node.addChild(names[1:], handlerFunc)
	}
}

func (n *ExactNode) isMatched(name string) bool {
	return n.name == "*" || n.name == name
}

func (n *ExactNode) findChild(name string) Node {
	if child := n.children[name]; child != nil {
		return child
	} else if child := n.children["*"]; child != nil {
		return child
	}
	return nil
}

func (n *ExactNode) findHandlerFunc(names []string) HandlerFunc {
	if matched := n.findChild(names[0]); matched != nil {
		if len(names) > 1 {
			if handlerFunc := matched.findHandlerFunc(names[1:]); handlerFunc != nil {
				return handlerFunc
			}
		}
		return matched.getHandlerFunc()
	}
	return n.getHandlerFunc()
}
