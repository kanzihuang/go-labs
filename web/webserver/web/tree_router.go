package web

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Node interface {
	childOf(name string) Node
	findHandlerFunc(path string, paths []string) HandlerFunc
	isMatched(name string) bool
	addChild(names []string, handlerFunc HandlerFunc)
	wrapHandlerFunc(name string, handlerFunc HandlerFunc) HandlerFunc
	getHandlerFunc() HandlerFunc
	setHandlerFunc(handlerFunc HandlerFunc)
	getName() string
	getChildren() map[string]Node
	isDynamic() bool
	getDynamicNode() Node
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

func NewNode(name string, handlerFunc HandlerFunc) Node {
	var sign uint8
	if len(name) > 0 {
		sign = name[0]
	}
	switch sign {
	case '~':
		return NewRegNode(name[1:], handlerFunc)
	case ':':
		return NewParamNode(name[1:], handlerFunc)
	case '*':
		return NewAnyNode("*", handlerFunc)
	default:
		return NewBaseNode(name, handlerFunc)
	}
}

type BaseNode struct {
	name        string
	handlerFunc HandlerFunc
	children    map[string]Node
	dynamicNode Node
	dynamic     bool
}

func (n *BaseNode) getDynamicNode() Node {
	return n.dynamicNode
}

func NewBaseNode(name string, handlerFunc HandlerFunc) Node {
	return &BaseNode{
		name:        name,
		handlerFunc: handlerFunc,
		children:    make(map[string]Node),
	}
}

func (n *BaseNode) isDynamic() bool {
	return n.dynamic
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

type RegNode struct {
	BaseNode
	validPath *regexp.Regexp
}

func NewRegNode(name string, handlerFunc HandlerFunc) Node {
	return &RegNode{
		BaseNode: BaseNode{
			name:        "~",
			handlerFunc: handlerFunc,
			children:    make(map[string]Node),
			dynamic:     true,
		},
		validPath: regexp.MustCompile(name),
	}
}

func (n *RegNode) isMatched(path string) bool {
	return n.validPath.MatchString(path)
}

type AnyNode struct {
	BaseNode
}

func NewAnyNode(name string, handlerFunc HandlerFunc) Node {
	return &AnyNode{
		BaseNode: BaseNode{
			name:        name,
			handlerFunc: handlerFunc,
			children:    make(map[string]Node),
			dynamic:     true,
		},
	}
}

func (n *AnyNode) isMatched(path string) bool {
	return true
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
			dynamic:     true,
		},
	}
}

func (n *ParamNode) isMatched(path string) bool {
	return true
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

func (n *BaseNode) wrapHandlerFunc(_ string, handlerFunc HandlerFunc) HandlerFunc {
	if handlerFunc != nil {
		return handlerFunc
	} else {
		return n.handlerFunc
	}
}

//func (n *BaseNode) addNode(name string, handlerFunc HandlerFunc) Node {
//	key := name
//	if strings.HasPrefix(key, ":") {
//		key = ":"
//	}
//	var node Node
//	if key == ":" {
//		node = NewParamNode(name[1:], handlerFunc)
//	} else {
//		node = NewBaseNode(name, handlerFunc)
//	}
//	n.children[key] = node
//	return node
//}

func (n *BaseNode) addChild(names []string, handlerFunc HandlerFunc) {
	name := names[0]
	child, ok := n.children[name]
	if !ok {
		var f HandlerFunc
		if len(names) <= 1 {
			f = handlerFunc
		}
		//child = n.addNode(name, f)
		child = NewNode(name, f)
		if child.isDynamic() {
			if n.dynamicNode != nil {
				if child.getName() == n.dynamicNode.getName() {
					child = n.dynamicNode
				} else {
					panic(fmt.Errorf("duplicate registered routes: %s/%s", n.getName(), child.getName()))
				}
			} else {
				n.dynamicNode = child
			}
		} else {
			n.children[child.getName()] = child
		}
		if !child.isDynamic() {
			n.children[child.getName()] = child
		} else if n.dynamicNode == nil {
			n.dynamicNode = child
		}
	}
	if len(names) > 1 {
		child.addChild(names[1:], handlerFunc)
	}
}

func (n *BaseNode) isMatched(path string) bool {
	return n.name == path
}

func (n *BaseNode) childOf(path string) Node {
	if child := n.children[path]; child != nil {
		return child
	} else if n.dynamicNode != nil && n.dynamicNode.isMatched(path) {
		return n.dynamicNode
	}
	//else if child := n.children[":"]; child != nil {
	//	return child
	//} else if child := n.children["*"]; child != nil {
	//	return child
	//}

	return nil
}

func (n *BaseNode) findHandlerFunc(path string, paths []string) HandlerFunc {
	childPath := paths[0]
	var handlerFunc HandlerFunc
	if matched := n.childOf(childPath); matched != nil {
		if len(paths) > 1 {
			handlerFunc = matched.findHandlerFunc(childPath, paths[1:])
		} else {
			handlerFunc = matched.wrapHandlerFunc(childPath, nil)
		}
	}
	return n.wrapHandlerFunc(path, handlerFunc)
}

func (n *BaseNode) panicIfExist(paths []string) {
	//todo
}

func (n *ParamNode) findHandlerFunc(path string, paths []string) HandlerFunc {
	return n.wrapHandlerFunc(path, n.BaseNode.findHandlerFunc(path, paths))
}
