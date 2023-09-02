package web

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	W        http.ResponseWriter
	R        *http.Request
	ParamMap map[string]string
}

func NewContext() *Context {
	return &Context{
		ParamMap: map[string]string{},
	}
}

func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	c.W = w
	c.R = r
	c.ParamMap = make(map[string]string)
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
	return nil
}