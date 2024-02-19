package simplerpc

type Health struct {
}

const ServiceHealth = "Health"

func (h *Health) ServiceName() string {
	return ServiceHealth
}

func NewHealth() *Health {
	return &Health{}
}

func (h *Health) Ping(request *PingReq) (*PingResp, error) {
	return &PingResp{
		Payload: request.Payload,
	}, nil
}
