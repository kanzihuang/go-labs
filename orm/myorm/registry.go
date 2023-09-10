package myorm

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

type field struct {
	columnName string
	fieldName  string
}

type model struct {
	tableName string
	fieldMap  map[string]field
}

type modelOpt func(model *model) error

type registry struct {
	models sync.Map
}

func newRegistry() *registry {
	return &registry{}
}

func (r *registry) get(typ reflect.Type) (*model, error) {
	if m, ok := r.models.Load(typ); ok {
		return m.(*model), nil
	}
	return r.register(typ)
}

func (r *registry) register(typ reflect.Type, opts ...modelOpt) (*model, error) {
	m, err := r.parseModel(typ)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(m); err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, m)
	return m, nil
}

func (r *registry) parseModel(typ reflect.Type) (*model, error) {
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("数据类型不是指向结构体的指针: " + typ.Name())
	}
	typ = typ.Elem()
	m := &model{}
	if t, ok := reflect.Zero(typ).Interface().(TableName); ok {
		m.tableName = t.TableName()
	}
	if len(m.tableName) == 0 {
		m.tableName = fmt.Sprintf("`%s`", underlineName(typ.Name()))
	}
	m.fieldMap = make(map[string]field, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).IsExported() == false {
			continue
		}
		fieldName := typ.Field(i).Name
		m.fieldMap[fieldName] = field{
			columnName: underlineName(fieldName),
			fieldName:  fieldName,
		}
	}
	return m, nil
}

var reUnderlineName = struct {
	re   *regexp.Regexp
	repl string
}{
	re:   regexp.MustCompile("([a-z0-9]+)"),
	repl: "${1}_",
}

func underlineName(name string) string {
	res := string(reUnderlineName.re.ReplaceAll([]byte(name), []byte(reUnderlineName.repl)))
	return strings.ToLower(strings.Trim(res, "_"))
}
