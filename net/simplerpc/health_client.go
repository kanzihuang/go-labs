package simplerpc

import "encoding/json"

type HealthClient struct {
	proxy Proxy
}

func NewHealthClient(proxy Proxy) *HealthClient {
	return &HealthClient{
		proxy: proxy,
	}
}

func (h *HealthClient) Ping(request *PingReq) (*PingResp, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req := &Request{
		ServiceName: "Health",
		MethodName:  "Ping",
		Data:        data,
	}
	resp := h.proxy.Call(req)
	if resp.Err != nil {
		return nil, resp.Err.err
	}
	pingResp := &PingResp{}
	err = json.Unmarshal(resp.Data, pingResp)
	if err != nil {
		return nil, err
	}
	if len(pingResp.Payload) != len(request.Payload) {
		return nil, ErrInvalidPayloadLength
	}
	return pingResp, nil
}
