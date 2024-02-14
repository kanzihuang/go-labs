package simplerpc

type Request struct {
	ServiceName string
	MethodName  string
	Data        []byte
}
type Response struct {
	Data []byte
	Err  *RespError
}

type Proxy interface {
	Call(req *Request) *Response
}
