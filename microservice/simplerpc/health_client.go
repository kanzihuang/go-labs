package simplerpc

func NewHealthClient() *HealthClient {
	return &HealthClient{}
}

type HealthClient struct {
	Ping func(request *PingReq) (*PingResp, error)
}

func (h *HealthClient) ServiceName() string {
	return ServiceHealth
}
