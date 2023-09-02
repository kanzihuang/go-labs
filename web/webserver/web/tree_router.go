package web

import (
	"net/http"
	"strings"
)

type Node interface {
	findChild(name string) Node
	findHandlerFunc(path string, paths []string) HandlerFunc
	isMatched(name string) bool
	addChild(names []string, handlerFunc HandlerFunc)
	wrapHandlerFunc(name string, handlerFunc HandlerFunc) HandlerFunc
	getHandlerFunc() HandlerFunc
	setHandlerFunc(handlerFunc HandlerFunc)
	getName() string
	getChildren() map[string]Node
}

var supportMethods = [4]string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
}

type RouterBasedOnTree struct {
	forest map[string]Node
}

func (r *RouterBasedOnTree) Route(method string, pattern string, handlerFunc HandlerFunc) {
	pattern = strings.TrimRight(pattern, "/")
	names := strings.Split(pattern, "/")
	root, ok := r.forest[method]
	if !ok {
		root = NewBaseNode("", nil)
		r.forest[method] = root
	}
	if len(names) == 1 {
		root.setHandlerFunc(handlerFunc)
	} else {
		root.addChild(names[1:], handlerFunc)
	}
}

func (r *RouterBasedOnTree) FindHandlerFunc(method string, path string) HandlerFunc {
	path = strings.TrimRight(path, "/")
	names := strings.Split(path, "/")
	root, ok := r.forest[method]
	if ok {
		if len(names) <= 1 {
			return root.wrapHandlerFunc("", nil)
		} else {
			return root.findHandlerFunc(names[0], names[1:])
		}
	} else {
		return nil
	}
}

func NewRouterBasedOnTree() *RouterBasedOnTree {
	return &RouterBasedOnTree{
		forest: make(map[string]Node, len(supportMethods)),
	}
}

type BaseNode struct {
	name        string
	handlerFunc HandlerFunc
	children    map[string]Node
}

func (n *BaseNode) getChildren() map[string]Node {
	return n.children
}

func (n *BaseNode) getHandlerFunc() HandlerFunc {
	return n.handlerFunc
}

func (n *BaseNode) getName() string {
	return n.name
}

func NewBaseNode(name string, handlerFunc HandlerFunc) Node {
	return &BaseNode{
		name:        name,
		handlerFunc: handlerFunc,
		children:    make(map[string]Node),
	}
}

type ParamNode struct {
	BaseNode
}

func NewParamNode(name string, handlerFunc HandlerFunc) Node {
	return &ParamNode{
		BaseNode: BaseNode{
			name:        name,
			handlerFunc: handlerFunc,
			children:    make(map[string]Node),
		},
	}
}

func (n *ParamNode) wrapHandlerFunc(path string, handlerFunc HandlerFunc) HandlerFunc {
	return func(c *Context) {
		c.ParamMap[n.name] = path
		if handlerFunc != nil {
			handlerFunc(c)
		} else {
			n.handlerFunc(c)
		}
	}
}

func (n *BaseNode) setHandlerFunc(handlerFunc HandlerFunc) {
	n.handlerFunc = handlerFunc
}

func (n *BaseNode) wrapHandlerFunc(path string, handlerFunc HandlerFunc) HandlerFunc {
	if handlerFunc != nil {
		return handlerFunc
	} else {
		return n.handlerFunc
	}
}

func (n *BaseNode) addNode(name string, handlerFunc HandlerFunc) Node {
	key := name
	if strings.HasPrefix(key, ":") {
		key = ":"
	}
	var node Node
	if key == ":" {
		node = NewParamNode(name[1:], handlerFunc)
	} else {
		node = NewBaseNode(name, handlerFunc)
	}
	n.children[key] = node
	return node
}

func (n *BaseNode) addChild(names []string, handlerFunc HandlerFunc) {
	name := names[0]
	node := n.findChild(name)
	if node == nil {
		var f HandlerFunc
		if len(names) <= 1 {
			f = handlerFunc
		}
		node = n.addNode(name, f)
	}
	if len(names) > 1 {
		node.addChild(names[1:], handlerFunc)
	}
}

func (n *BaseNode) isMatched(name string) bool {
	return n.name == "*" || n.name == name
}

func (n *BaseNode) findChild(name string) Node {
	if child := n.children[name]; child != nil {
		return child
	} else if child := n.children[":"]; child != nil {
		return child
	} else if child := n.children["*"]; child != nil {
		return child
	}
	return nil
}

func (n *BaseNode) findHandlerFunc(path string, paths []string) HandlerFunc {
	childPath := paths[0]
	var handlerFunc HandlerFunc
	if matched := n.findChild(childPath); matched != nil {
		if len(paths) > 1 {
			handlerFunc = matched.findHandlerFunc(childPath, paths[1:])
		} else {
			handlerFunc = matched.wrapHandlerFunc(childPath, nil)
		}
	}
	return n.wrapHandlerFunc(path, handlerFunc)
}

func (n *ParamNode) findHandlerFunc(path string, paths []string) HandlerFunc {
	return n.wrapHandlerFunc(path, n.BaseNode.findHandlerFunc(path, paths))
}
