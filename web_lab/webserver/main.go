package main

import (
	"fmt"
	"go-labs/web_lab/webserver/web"
	"io"
	"log"
	"net/http"
)

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

func main() {
	server := web.NewServer()
	server.Route(http.MethodGet, "/", home)
	server.Route(http.MethodGet, "/header", header)
	server.Route(http.MethodPost, "/form", form)
	server.Route(http.MethodPost, "/body", body)
	server.Route(http.MethodPost, "/user/signup", signUp)
	server.Route(http.MethodGet, "/user/*", user)
	log.Fatal(server.Start(":8080"))
}

func user(c *web.Context) {
	_ = c.WriteJson(http.StatusOK, &commonResponse{
		Data: c.R.URL.Path,
	})
}

func signUp(ctx *web.Context) {
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

func body(ctx *web.Context) {
	b, err := io.ReadAll(ctx.R.Body)
	if err != nil {
		fmt.Fprintln(ctx.W, "read all err: ", err)
		return
	}
	ctx.W.Write(b)
}

func form(ctx *web.Context) {
	err := ctx.R.ParseForm()
	if err != nil {
		fmt.Fprintf(ctx.W, "parse form error: %v\n", err)
		return
	}
	fmt.Fprintf(ctx.W, "form: %v\n", ctx.R.Form)
}

func header(ctx *web.Context) {
	fmt.Fprintf(ctx.W, "header: %v\n", ctx.R.Header)
}

func home(ctx *web.Context) {
	ctx.W.WriteHeader(http.StatusOK)
	fmt.Fprintf(ctx.W, "Hi there, I love %s!\n", ctx.R.URL.Path[1:])
}
