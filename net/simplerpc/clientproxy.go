package simplerpc

type ClientProxy struct {
	client *Client
}

func (s *ClientProxy) Call(req *Request) *Response {
	return s.client.Call(req)
}
