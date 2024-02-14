package simplerpc

type Health struct {
}

func NewHealth() *Health {
	return &Health{}
}

func (h *Health) Ping(payload []byte) ([]byte, error) {
	return payload, nil
}
