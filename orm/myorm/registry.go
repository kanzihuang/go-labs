package myorm

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

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

func (r *registry) parseTagColumn(tag reflect.StructTag) (string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return "", nil
	}
	items := strings.Split(string(ormTag), ",")
	for _, item := range items {
		peer := strings.Split(item, "=")
		if len(peer) != 2 {
			continue
		}
		if peer[0] == "column" {
			return peer[1], nil
		}
	}
	return "", nil
}

func (r *registry) parseModel(typ reflect.Type) (*model, error) {
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("数据类型不是指向结构体的指针: " + typ.Name())
	}
	typ = typ.Elem()
	m := newModel()
	if t, ok := reflect.Zero(typ).Interface().(TableName); ok {
		m.tableName = t.TableName()
	}
	if len(m.tableName) == 0 {
		m.tableName = fmt.Sprintf("`%s`", underlineName(typ.Name()))
	}

	m.fieldMap = make(map[string]field, typ.NumField())
	m.columnMap = make(map[string]field, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		tf := typ.Field(i)
		if tf.IsExported() == false {
			continue
		}
		colName, err := r.parseTagColumn(tf.Tag)
		if err != nil {
			return nil, err
		}
		if len(colName) == 0 {
			colName = underlineName(tf.Name)
		}
		mf := field{
			columnName: colName,
			fieldName:  tf.Name,
			typ:        tf.Type,
		}
		m.fieldMap[tf.Name] = mf
		m.columnMap[colName] = mf
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
