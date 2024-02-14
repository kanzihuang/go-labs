package simplerpc

type HealthClient struct {
	proxy Proxy
}

func NewHealthClient(proxy Proxy) *HealthClient {
	return &HealthClient{
		proxy: proxy,
	}
}

func (h *HealthClient) Ping(payload []byte) ([]byte, error) {
	req := &Request{
		ServiceName: "Health",
		MethodName:  "Ping",
		Arg:         payload,
	}
	resp := h.proxy.Call(req)
	if resp.Err != nil {
		return nil, resp.Err.err
	}
	if len(resp.Resp) != len(payload) {
		return nil, ErrInvalidPayloadLength
	}
	return resp.Resp, nil
}
