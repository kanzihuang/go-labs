package simplerpc

import (
	"encoding/json"
	"errors"
)

var (
	ErrInvalidPayloadType   = errors.New("simplerpc: 无效的载荷类型")
	ErrInvalidPayloadLength = errors.New("simplerpc: 无效的载荷长度")
	ErrInvalidServiceName   = errors.New("simplerpc: 无效的服务名称")
	ErrInvalidMethodName    = errors.New("simplerpc: 无效的方法名称")
	ErrMultipleServiceName  = errors.New("simplerpc: 重复的服务名称")
)

func NewError(err error) *RespError {
	return &RespError{
		err: err,
	}
}

type RespError struct {
	err error
}

func (e *RespError) UnmarshalJSON(bytes []byte) error {
	var str string
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		return err
	}
	e.err = nil
	if len(str) > 0 {
		e.err = errors.New(str)
	}
	return nil
}

func (e *RespError) MarshalJSON() ([]byte, error) {
	err := ""
	if e.err != nil {
		err = e.err.Error()
	}
	return json.Marshal(err)
}

func (e *RespError) Error() string {
	if e != nil && e.err != nil {
		return e.err.Error()
	}
	return ""
}
