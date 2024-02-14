package simplerpc

const (
	MethodArgsNum    = 1
	MethodResultsNum = 2
)

type PingReq struct {
	Payload []byte
}

type PingResp struct {
	Payload []byte
}

type Service interface {
	ServiceName() string
}
