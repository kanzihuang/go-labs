package simplerpc

import (
	"encoding/json"
	"reflect"
)

type ServerProxy struct {
	services map[string]Service
}

func NewServerProxy() *ServerProxy {
	return &ServerProxy{
		services: make(map[string]Service, 16),
	}
}

func (s *ServerProxy) Register(svc Service) error {
	name := svc.ServiceName()
	_, ok := s.services[name]
	if ok {
		return ErrMultipleServiceName
	}
	s.services[name] = svc
	return nil
}

func (s *ServerProxy) call(svc Service, methodName string, data []byte) (any, error) {
	svcVal := reflect.ValueOf(svc)
	method := svcVal.MethodByName(methodName)
	if method.IsZero() {
		return nil, ErrInvalidMethodName
	}
	methodType := method.Type()
	reqType := methodType.In(0)
	reqVal := reflect.New(reqType.Elem())
	req := reqVal.Interface()
	err := json.Unmarshal(data, req)
	if err != nil {
		return nil, err
	}
	results := method.Call([]reflect.Value{reqVal})
	if len(results) != 2 {
		return nil, ErrInvalidMethodResultsNum
	}
	if res1 := results[1].Interface(); res1 != nil {
		err, ok := res1.(error)
		if !ok {
			return nil, ErrInvalidMethodResultsType
		}
		return nil, err
	}
	resp := results[0].Interface()
	return resp, nil
}

func (s *ServerProxy) Call(req *ProxyReq) *ProxyResp {
	svc, ok := s.services[req.ServiceName]
	if !ok {
		return &ProxyResp{
			Data: nil,
			Err:  NewError(ErrInvalidServiceName),
		}
	}
	resp, err := s.call(svc, req.MethodName, req.Data)
	if err != nil {
		return &ProxyResp{
			Data: nil,
			Err:  NewError(err),
		}
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return &ProxyResp{
			Data: nil,
			Err:  NewError(err),
		}
	}
	return &ProxyResp{
		Data: data,
		Err:  nil,
	}
}
