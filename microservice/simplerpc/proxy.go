package simplerpc

type ProxyReq struct {
	ServiceName string
	MethodName  string
	Data        []byte
}
type ProxyResp struct {
	Data []byte
	Err  *RespError
}

type Proxy interface {
	Call(req *ProxyReq) *ProxyResp
	Register(svc Service) error
}
