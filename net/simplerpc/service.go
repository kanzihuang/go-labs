package simplerpc

type Service interface {
	Ping(payload []byte) ([]byte, error)
}
