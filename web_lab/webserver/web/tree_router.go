package web

import (
	"log"
	"net/http"
	"strings"
)

type Node interface {
	findChild(name string) Node
	findHandlerFunc(path string, paths []string) HandlerFunc
	isMatched(name string) bool
	addChild(names []string, handlerFunc HandlerFunc)
	handlerFuncBuilder(name string, handlerFunc HandlerFunc) HandlerFunc
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
	pattern = strings.TrimRight(pattern, "/")
	names := strings.Split(pattern, "/")
	root, ok := r.forest[method]
	if ok {
		if len(names) == 1 {
			root.setHandlerFunc(handlerFunc)
		} else {
			root.addChild(names[1:], handlerFunc)
		}
	} else {
		log.Panic("Not supported method", method)
	}
}

func (r *RouterBasedOnTree) FindHandlerFunc(method string, path string) HandlerFunc {
	path = strings.TrimRight(path, "/")
	names := strings.Split(path, "/")
	root, ok := r.forest[method]
	if ok {
		if len(names) == 0 {
			return root.handlerFuncBuilder("", nil)
		} else {
			return root.findHandlerFunc(names[0], names[1:])
		}
	} else {
		return nil
	}
}

func NewRouterBasedOnTree() Router {
	forest := make(map[string]Node, len(supportMethods))
	for _, method := range supportMethods {
		forest[method] = NewBaseNode("", nil)
	}
	return &RouterBasedOnTree{
		forest: forest,
	}
}

type BaseNode struct {
	name        string
	handlerFunc HandlerFunc
	children    map[string]Node
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

func (n *ParamNode) handlerFuncBuilder(name string, handlerFunc HandlerFunc) HandlerFunc {
	return func(c *Context) {
		c.ParamMap[n.name] = name
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

func (n *BaseNode) handlerFuncBuilder(string, HandlerFunc) HandlerFunc {
	return n.handlerFunc
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
	name := paths[0]
	if matched := n.findChild(name); matched != nil {
		if len(paths) > 1 {
			if handlerFunc := matched.findHandlerFunc(name, paths[1:]); handlerFunc != nil {
				return handlerFunc
			}
		}
		return matched.handlerFuncBuilder(name, nil)
	}
	return n.handlerFuncBuilder(path, nil)
}

func (n *ParamNode) findHandlerFunc(path string, paths []string) HandlerFunc {
	return n.handlerFuncBuilder(path, n.BaseNode.findHandlerFunc(path, paths))
}
