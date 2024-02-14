package simplerpc

import (
	"encoding/json"
	"reflect"
)

func NewClientProxy(client *Client) *ClientProxy {
	return &ClientProxy{
		client: client,
	}
}

type ClientProxy struct {
	client *Client
}

func (s *ClientProxy) Call(req *ProxyReq) *ProxyResp {
	return s.client.Call(req)
}

func (s *ClientProxy) Register(svc Service) error {
	return s.setFuncFields(svc)
}

func (s *ClientProxy) setFuncFields(svc Service) error {
	svcVal := reflect.ValueOf(svc).Elem()
	svcTyp := svcVal.Type()
	for i := 0; i < svcVal.NumField(); i++ {
		fieldInfo := svcTyp.Field(i)
		fieldVal := svcVal.Field(i)
		if fieldVal.Kind() != reflect.Func {
			continue
		}
		fnVal, err := s.makeFunc(fieldInfo.Name, fieldInfo.Type)
		if err != nil {
			return err
		}
		fieldVal.Set(fnVal)
	}
	return nil
}

func (s *ClientProxy) makeFunc(methodName string, fieldTyp reflect.Type) (reflect.Value, error) {
	//fieldTyp := fieldVal.Type()
	if fieldTyp.NumIn() != MethodArgsNum {
		return reflect.Value{}, ErrInvalidMethodArgsNum
	}
	if fieldTyp.NumOut() != MethodResultsNum {
		return reflect.Value{}, ErrInvalidMethodResultsNum
	}
	fnVal := reflect.MakeFunc(fieldTyp, func(args []reflect.Value) (results []reflect.Value) {
		reqZero := reflect.Zero(fieldTyp.In(0))
		if len(args) != fieldTyp.NumIn() {
			return []reflect.Value{reqZero, reflect.ValueOf(ErrInvalidMethodArgsNum)}
		}
		req := args[0].Interface()
		resp, err := s.call(methodName, req, fieldTyp.Out(0))
		if err != nil {
			return []reflect.Value{reqZero, reflect.ValueOf(err)}
		}
		return []reflect.Value{
			reflect.ValueOf(resp),
			reflect.Zero(reflect.TypeOf((*error)(nil)).Elem()),
		}
	})
	return fnVal, nil
}

func (s *ClientProxy) call(methodName string, request any, respTyp reflect.Type) (any, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req := &ProxyReq{
		ServiceName: ServiceHealth,
		MethodName:  methodName,
		Data:        data,
	}
	resp := s.Call(req)
	if resp.Err != nil {
		return nil, resp.Err.err
	}
	result := reflect.New(respTyp.Elem()).Interface()
	err = json.Unmarshal(resp.Data, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
