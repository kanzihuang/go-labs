package simplerpc

type ServerProxy struct {
	services map[string]Service
}

func NewServerProxy() *ServerProxy {
	return &ServerProxy{
		services: make(map[string]Service, 16),
	}
}

func (s *ServerProxy) Register(name string, svc Service) error {
	_, ok := s.services[name]
	if ok {
		return ErrMultipleServiceName
	}
	s.services[name] = svc
	return nil
}

func (s *ServerProxy) Call(req *Request) *Response {
	svc, ok := s.services[req.ServiceName]
	if !ok {
		return &Response{
			Resp: nil,
			Err:  NewError(ErrInvalidServiceName),
		}
	}
	respPayload, err := svc.Ping(req.Arg)
	if err != nil {
		return &Response{
			Resp: nil,
			Err:  NewError(err),
		}
	}
	return &Response{
		Resp: respPayload,
		Err:  nil,
	}
}
