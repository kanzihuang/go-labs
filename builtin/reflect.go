package builtin

import (
	"errors"
	"reflect"
)

func getFieldPointer(obj any, fieldName string) (any, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Pointer {
		return nil, errors.New("value kind is not Pointer: " + val.Kind().String())
	}
	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return nil, errors.New("elem kind is not Struct: " + val.Kind().String())
	}
	field := elem.FieldByName(fieldName)
	if field.Kind() == reflect.Invalid {
		return nil, errors.New("field name is not found: " + fieldName)
	}
	return reflect.NewAt(field.Type(), field.Addr().UnsafePointer()).Interface(), nil
}
