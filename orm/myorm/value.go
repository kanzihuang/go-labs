package myorm

import (
	"database/sql"
	"reflect"
)

type valueCreator func(val any, meta *model) value
type value interface {
	SetColumns(rows *sql.Rows) error
}

type reflectValue struct {
	val  reflect.Value
	meta *model
}

func newReflectValue(val any, meta *model) value {
	return &reflectValue{
		val:  reflect.ValueOf(val).Elem(),
		meta: meta,
	}
}

func (val *reflectValue) SetColumns(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	lenCols := len(columns)
	values := make([]any, lenCols)
	elemValues := make([]reflect.Value, lenCols)
	for i, column := range columns {
		typ := val.meta.columnMap[column].typ
		v := reflect.New(typ)
		values[i] = v.Interface()
		elemValues[i] = v.Elem()
	}

	if err := rows.Scan(values...); err != nil {
		return err
	}
	for i, column := range columns {
		fld := val.val.FieldByName(val.meta.columnMap[column].fieldName)
		fld.Set(elemValues[i])
	}
	return nil
}

type unsafeValue struct {
	val  reflect.Value
	meta *model
}

func newUnsafeValue(val any, meta *model) value {
	return &unsafeValue{
		val:  reflect.ValueOf(val).Elem(),
		meta: meta,
	}
}

func (val *unsafeValue) SetColumns(rows *sql.Rows) error {
	return nil
}
