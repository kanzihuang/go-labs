package myorm

import (
	"reflect"
)

type field struct {
	columnName string
	fieldName  string
	typ        reflect.Type
}

type model struct {
	tableName string
	fieldMap  map[string]field
	columnMap map[string]field
	creator   valueCreator
}

func newModel() *model {
	return &model{}
}
