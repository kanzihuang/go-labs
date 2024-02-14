package simplerpc

type PingReq struct {
	Payload []byte
}

type PingResp struct {
	Payload []byte
}

type Service interface {
	Ping(request *PingReq) (*PingResp, error)
}
