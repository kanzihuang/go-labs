package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Routable interface {
	Route(method string, pattern string, handlerFunc func(c *Context))
}

type Server = interface {
	Routable
	Start(address string) error
}

type Handler interface {
	http.Handler
	Routable
}

type sdkHttpServer struct {
	Name    string
	handler Handler
}

func (s *sdkHttpServer) Route(method, pattern string, handlerFunc func(c *Context)) {
	s.handler.Route(method, pattern, handlerFunc)
}

func (s *sdkHttpServer) Start(address string) error {
	http.Handle("/", s.handler)
	return http.ListenAndServe(address, nil)
}

func NewServer() Server {
	return &sdkHttpServer{
		Name:    "sdkHttpServer",
		handler: NewHandlerBasedOnMap(),
	}
}

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
	}
}
func (c *Context) ReadJson(data interface{}) error {
	body, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}

func (c *Context) WriteJson(status int, data interface{}) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.W.Write(bs)
	if err != nil {
		return err
	}
	c.W.WriteHeader(status)
	return nil
}

type signUpReq struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int    `json:"biz_code"`
	Msg     string `json:"message"`
	Data    string `json:"data"`
}

type HandlerBasedOnMap struct {
	handlers map[string]func(c *Context)
}

var _ Handler = (*HandlerBasedOnMap)(nil)

func NewHandlerBasedOnMap() Handler {
	return &HandlerBasedOnMap{
		handlers: map[string]func(c *Context){},
	}
}

func (h *HandlerBasedOnMap) Route(method string, pattern string, handlerFunc func(c *Context)) {
	key := h.key(method, pattern)
	h.handlers[key] = handlerFunc
}

func (h *HandlerBasedOnMap) key(method string, path string) string {
	return fmt.Sprintf("%s#%s", method, path)
}

func (h *HandlerBasedOnMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := h.key(r.Method, r.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		c := &Context{W: w, R: r}
		handler(c)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not any router matched"))
	}
}

func main() {
	server := NewServer()
	server.Route(http.MethodGet, "/", home)
	server.Route(http.MethodGet, "/header", header)
	server.Route(http.MethodPost, "/form", form)
	server.Route(http.MethodPost, "/body", body)
	server.Route(http.MethodPost, "/signup", signUp)
	log.Fatal(server.Start(":8080"))
}

func signUp(ctx *Context) {
	s := &signUpReq{}
	err := ctx.ReadJson(s)
	if err != nil {
		_ = ctx.WriteJson(http.StatusOK, &commonResponse{
			BizCode: http.StatusBadRequest,
			Msg:     fmt.Sprintf("invalid request: %v\n", err),
		})
		return
	}
	_ = ctx.WriteJson(http.StatusOK, &commonResponse{
		Data: "123",
	})
}

func body(ctx *Context) {
	b, err := io.ReadAll(ctx.R.Body)
	if err != nil {
		fmt.Fprintln(ctx.W, "read all err: ", err)
		return
	}
	ctx.W.Write(b)
}

func form(ctx *Context) {
	err := ctx.R.ParseForm()
	if err != nil {
		fmt.Fprintf(ctx.W, "parse form error: %v\n", err)
		return
	}
	fmt.Fprintf(ctx.W, "form: %v\n", ctx.R.Form)
}

func header(ctx *Context) {
	fmt.Fprintf(ctx.W, "header: %v\n", ctx.R.Header)
}

func home(ctx *Context) {
	ctx.W.WriteHeader(http.StatusOK)
	fmt.Fprintf(ctx.W, "Hi there, I love %s!\n", ctx.R.URL.Path[1:])
}
