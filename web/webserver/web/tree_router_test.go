package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"testing"
)

func nameOfHandlerFunc(fn HandlerFunc) string {
	if fn == nil {
		return "nil"
	}
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func (r *RouterBasedOnTree) equal(dst Router) (string, bool) {
	if reflect.TypeOf(r) != reflect.TypeOf(dst) {
		return "routers are not of equal type", false
	}
	if msg, ok := equalNodeMap(r.forest, dst.(*RouterBasedOnTree).forest); !ok {
		return "inconsistent router, " + msg, ok
	}
	return "", true
}

func equalNodeMap(x map[string]Node, y map[string]Node) (string, bool) {
	if len(x) != len(y) {
		return "node maps are not of equal length", false
	}
	for k, a := range x {
		b := y[k]
		if msg, ok := equalNode(a, b); ok != true {
			return "inconsistent node, " + msg, ok
		}
	}
	return "", true
}

func equalNode(x, y Node) (string, bool) {
	if x == nil || y == nil {
		return "not all nodes are nil", x == nil && y == nil
	}
	if reflect.TypeOf(x) != reflect.TypeOf(y) {
		return "inconsistent node types", false
	}
	if x.getName() != y.getName() {
		return "inconsistent node names", false
	}
	if nameOfHandlerFunc(x.getHandlerFunc()) != nameOfHandlerFunc(y.getHandlerFunc()) {
		return "inconsistent node handler", false
	}
	if msg, ok := equalNodeMap(x.getChildren(), y.getChildren()); !ok {
		return "inconsistent node maps, " + msg, ok
	}
	if msg, ok := equalNode(x.getDynamicNode(), y.getDynamicNode()); !ok {
		return "inconsistent dynamic node, " + msg, ok
	}
	return "", true
}

func handleRoot(*Context)        {}
func handleLogin(*Context)       {}
func handleConfig(*Context)      {}
func handleConfigPort(*Context)  {}
func handleStaticImage(*Context) {}
func handleRegexp(*Context)      {}
func handleRegexpInfo(*Context)  {}
func handleStaticAny(*Context)   {}
func handleStaticSize(*Context)  {}
func handleOrder(c *Context) {
	c.ParamMap["handlerFunc"] = "handleOrder"
}
func handleOrderStatus(c *Context) {
	c.ParamMap["handlerFunc"] = "handleOrderStatus"
}

type testCases []*struct {
	method      string            // if method is empty, assign http.MethodGet to method
	pattern     string            // if pattern is empty, don't registry route
	path        string            // if path is empty, assign pattern to path
	params      map[string]string // for param route
	handlerFunc HandlerFunc
}

func newTestCases() testCases {
	tcs := testCases{
		// fix route
		{
			pattern:     "/",
			handlerFunc: handleRoot,
		},
		{
			pattern:     "/login",
			handlerFunc: handleLogin,
		},
		{
			pattern:     "/config/",
			handlerFunc: handleConfig,
		},
		{
			pattern:     "/config/port",
			handlerFunc: handleConfigPort,
		},
		{
			pattern:     "/static/img",
			handlerFunc: handleStaticImage,
		},
		{
			path:        "/no_routing",
			handlerFunc: handleRoot,
		},
		{
			path:        "/login/no_routing",
			handlerFunc: handleLogin,
		},
		{
			path:        "/no_routing/no_routing",
			handlerFunc: handleRoot,
		},
		// regexp route
		{
			pattern:     "/regexp_parent/~^reg[a-z]+",
			path:        "/regexp_parent/regexp",
			handlerFunc: handleRegexp,
		},
		{
			path:        "/regexp_parent/reg-exp/child",
			handlerFunc: handleRoot,
		},
		{
			path:        "/regexp_parent/regexp/child",
			handlerFunc: handleRegexp,
		},
		{
			pattern:     "/regexp_parent/~^reg[a-z]+/info",
			path:        "/regexp_parent/regexp/info",
			handlerFunc: handleRegexpInfo,
		},
		// param route
		{
			pattern: "/order/:id",
			path:    "/order/3721",
			params: map[string]string{
				"id":          "3721",
				"handlerFunc": "handleOrder",
			},
			handlerFunc: handleOrder,
		},
		{
			pattern: "/order/:id/status",
			path:    "/order/3721/status",
			params: map[string]string{
				"id":          "3721",
				"handlerFunc": "handleOrderStatus",
			},
			handlerFunc: handleOrderStatus,
		},
		// any route
		{
			pattern:     "/static/*",
			path:        "/static/any",
			handlerFunc: handleStaticAny,
		},
		{
			path:        "/static/any/any/1.jpg",
			handlerFunc: handleStaticAny,
		},
		{
			pattern:     "/static/*/size",
			path:        "/static/any/size",
			handlerFunc: handleStaticSize,
		},
	}
	for _, tc := range tcs {
		if tc.method == "" {
			tc.method = http.MethodGet
		}
		if tc.path == "" {
			tc.path = tc.pattern
		}
	}
	return tcs
}

func (tcs testCases) registry(router *RouterBasedOnTree) {
	for _, tc := range tcs {
		//if !strings.HasPrefix(tc.pattern, "/regexp_parent/~^reg[a-z]+") {
		//	//todo comment this code
		//	continue
		//}
		if tc.pattern != "" {
			router.Route(tc.method, tc.pattern, tc.handlerFunc)
		}
	}
}

func TestRouterBasedOnTree_Route(t *testing.T) {
	r := NewRouterBasedOnTree()
	tcs := newTestCases()
	tcs.registry(r)

	wanted := RouterBasedOnTree{
		forest: map[string]Node{
			http.MethodGet: &BaseNode{
				name:        "",
				handlerFunc: handleRoot,
				children: map[string]Node{
					"login": &BaseNode{
						name:        "login",
						handlerFunc: handleLogin,
					},
					"config": &BaseNode{
						name:        "config",
						handlerFunc: handleConfig,
						children: map[string]Node{
							"port": &BaseNode{
								name:        "port",
								handlerFunc: handleConfigPort,
							},
						},
					},
					"static": &BaseNode{
						name:        "static",
						handlerFunc: nil,
						children: map[string]Node{
							"img": &BaseNode{
								name:        "img",
								handlerFunc: handleStaticImage,
							},
						},
						dynamicNode: &AnyNode{
							BaseNode{
								name:        "*",
								handlerFunc: handleStaticAny,
								children: map[string]Node{
									"size": &BaseNode{
										name:        "size",
										handlerFunc: handleStaticSize,
									},
								},
							},
						},
					},
					"order": &BaseNode{
						name:        "order",
						handlerFunc: nil,
						dynamicNode: &ParamNode{
							BaseNode{
								name:        "id",
								handlerFunc: handleOrder,
								children: map[string]Node{
									"status": &BaseNode{
										name:        "status",
										handlerFunc: handleOrderStatus,
									},
								},
							},
						},
					},
					"regexp_parent": &BaseNode{
						name:        "regexp_parent",
						handlerFunc: nil,
						dynamicNode: &RegNode{
							BaseNode: BaseNode{
								name:        "~",
								handlerFunc: handleRegexp,
								children: map[string]Node{
									"info": &BaseNode{
										name:        "info",
										handlerFunc: handleRegexpInfo,
									},
								},
							},
							validPath: regexp.MustCompile("^reg[a-z]+"),
						},
					},
				},
			},
		},
	}

	if msg, ok := wanted.equal(r); !ok {
		t.Errorf("error: %s", msg)
	}
}

func TestRouterBasedOnTree_FindHandlerFunc(t *testing.T) {
	r := NewRouterBasedOnTree()
	tcs := newTestCases()
	tcs.registry(r)
	for _, tc := range tcs {
		//if tc.path != "/regexp_parent/regexp/info" {
		//	continue
		//}
		f := r.FindHandlerFunc(tc.method, tc.path)
		assertEqualHandlerFunc(t, tc.handlerFunc, f, tc.path, tc.params)
	}
}

func assertEqualHandlerFunc(t *testing.T, expected, actual HandlerFunc, path string, params map[string]string) {
	if len(params) == 0 {
		assert.Equal(t, nameOfHandlerFunc(expected), nameOfHandlerFunc(actual),
			fmt.Sprintf("unable to find correct handlerFunc for path '%s'", path))
	} else {
		c := NewContext()
		actual(c)
		assert.Equal(t, params, c.ParamMap, "path: %s", path)
	}
}
