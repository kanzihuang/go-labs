package simplerpc

type Request struct {
	ServiceName string
	MethodName  string
	Arg         []byte
}
type Response struct {
	Resp []byte
	Err  *ResponseError
}

type Proxy interface {
	Call(req *Request) *Response
}
