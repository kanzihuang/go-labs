package myorm

import (
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
	columnName string
}

func C(columnName string) Column {
	return Column{columnName: columnName}
}

func (c Column) Build() (string, error) {
	return "`" + c.columnName + "`", nil
}

type Value struct {
}

func (v Value) Build() (string, error) {
	return "?", nil
}

type Predicate struct {
	left    Expression
	operate Operate
	right   Expression
}

func (p Predicate) Build() (string, error) {
	sb := strings.Builder{}
	switch p.operate {
	case OperateEqual, OperateAnd, OperateOr:
		if res, err := p.left.Build(); err == nil {
			sb.WriteString(res)
		} else {
			return "", err
		}
		sb.WriteByte(' ')
		sb.WriteString(string(p.operate))
		sb.WriteByte(' ')
		if res, err := p.right.Build(); err == nil {
			sb.WriteString(res)
		} else {
			return "", err
		}
	case OperateNot:
		sb.WriteByte(' ')
		sb.WriteString(string(p.operate))
		sb.WriteByte(' ')
		if res, err := p.right.Build(); err == nil {
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
	tableName string
	where     string
	args      []any
}

func (s Selector[T]) From(tableName string) QueryBuilder {
	s.tableName = tableName
	return s
}

func (s Selector[T]) Build() (*Query, error) {
	sb := strings.Builder{}

	sb.WriteString("SELECT * FROM ")

	if len(s.tableName) > 0 {
		sb.WriteString(s.tableName)
	} else {
		sb.WriteByte('`')
		sb.WriteString(reflect.TypeOf(*new(T)).Name())
		sb.WriteByte('`')
	}

	if len(s.where) > 0 {
		sb.WriteString(s.where)
	}

	sb.WriteString(";")
	return &Query{
		SQL:    sb.String(),
		Params: s.args,
	}, nil
}

func (s Selector[T]) Where(expr Expression, args []any) QueryBuilder {
	where, _ := expr.Build()
	length := len(where)
	if length > 0 {
		if where[0] == '(' && where[length-1] == ')' {
			where = where[1 : length-1]
		}
		s.where = " WHERE " + strings.TrimSpace(where)
		s.args = args
	} else {
		s.where = ""
		s.args = nil
	}
	return s
}
