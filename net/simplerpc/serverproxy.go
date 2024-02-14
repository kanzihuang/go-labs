package simplerpc

import "encoding/json"

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
			Data: nil,
			Err:  NewError(ErrInvalidServiceName),
		}
	}
	pingReq := &PingReq{}
	err := json.Unmarshal(req.Data, pingReq)
	if err != nil {
		return &Response{
			Data: nil,
			Err:  NewError(err),
		}
	}
	pingResp, err := svc.Ping(pingReq)
	if err != nil {
		return &Response{
			Data: nil,
			Err:  NewError(err),
		}
	}
	data, err := json.Marshal(pingResp)
	if err != nil {
		return &Response{
			Data: nil,
			Err:  NewError(err),
		}
	}
	return &Response{
		Data: data,
		Err:  nil,
	}
}
