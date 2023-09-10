package myorm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Operate string

const (
	OperateEqual Operate = "="
	OperateAnd   Operate = "AND"
	OperateOr    Operate = "OR"
	OperateNot   Operate = "NOT"
)

type Column struct {
	fieldName string
}

func C(fieldName string) Column {
	return Column{fieldName: fieldName}
}

func (c Column) Build(model *model) (string, error) {
	f, ok := model.fieldMap[c.fieldName]
	if ok != true {
		return "", errors.New("invalid column: " + c.fieldName)
	}
	return "`" + f.columnName + "`", nil
}

type Value struct {
}

func (v Value) Build(*model) (string, error) {
	return "?", nil
}

type Predicate struct {
	left    Expression
	operate Operate
	right   Expression
}

func (p Predicate) Build(model *model) (string, error) {
	sb := strings.Builder{}
	switch p.operate {
	case OperateEqual, OperateAnd, OperateOr:
		if res, err := p.left.Build(model); err == nil {
			sb.WriteString(res)
		} else {
			return "", err
		}
		sb.WriteByte(' ')
		sb.WriteString(string(p.operate))
		sb.WriteByte(' ')
		if res, err := p.right.Build(model); err == nil {
			sb.WriteString(res)
		} else {
			return "", err
		}
	case OperateNot:
		sb.WriteByte(' ')
		sb.WriteString(string(p.operate))
		sb.WriteByte(' ')
		if res, err := p.right.Build(model); err == nil {
			sb.WriteString(res)
		} else {
			return "", err
		}
	default:
		return "", fmt.Errorf("orm: 不支持的操作符: %v", p.operate)
	}

	result := strings.TrimSpace(sb.String())
	if _, ok := p.left.(Predicate); ok {
		result = "(" + result + ")"
	}
	return result, nil
}

type Selector[T any] struct {
	db       *DB
	typePtrT reflect.Type
	where    Expression
	args     []any
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db:       db,
		typePtrT: reflect.TypeOf(new(T)),
	}
}

func (s Selector[T]) buildEnd(sb *strings.Builder) {
	sb.WriteString(";")
}

func (s Selector[T]) buildSelect(sb *strings.Builder) {
	sb.WriteString("SELECT * FROM ")
}

func (s Selector[T]) buildFrom(sb *strings.Builder) error {
	m, err := s.db.registry.get(s.typePtrT)
	if err != nil {
		return err
	}
	sb.WriteString(m.tableName)
	return nil
}

func (s Selector[T]) buildWhere(sb *strings.Builder) error {
	if s.where == nil {
		return nil
	}
	m, err := s.db.registry.get(s.typePtrT)
	where, err := s.where.Build(m)
	if err != nil {
		return err
	}
	length := len(where)
	if length > 0 {
		if where[0] == '(' && where[length-1] == ')' {
			where = where[1 : length-1]
		}
		where = " WHERE " + strings.TrimSpace(where)
		sb.WriteString(where)
	}
	return nil
}

func (s Selector[T]) Build() (*Query, error) {
	sb := strings.Builder{}
	s.buildSelect(&sb)
	if err := s.buildFrom(&sb); err != nil {
		return nil, err
	}
	if err := s.buildWhere(&sb); err != nil {
		return nil, err
	}
	s.buildEnd(&sb)
	return &Query{
		SQL:    sb.String(),
		Params: s.args,
	}, nil
}

func (s Selector[T]) Where(expr Expression, args []any) QueryBuilder {
	s.where = expr
	s.args = args
	return s
}
