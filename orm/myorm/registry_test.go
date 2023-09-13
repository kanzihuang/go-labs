package myorm

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_underlineName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "AaaBbCc",
			want: "aaa_bb_cc",
		},
		{
			name: "AAaB",
			want: "aaa_b",
		},
		{
			name: "Aa5B",
			want: "aa5_b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, underlineName(tt.name), "underlineName(%v)", tt.name)
		})
	}
}

func Test_parseModel(t *testing.T) {
	type testModel struct {
		Id          uint64
		FirstName   string
		nonExported string
	}
	testCases := []struct {
		name      string
		val       any
		wantErr   error
		wantModel *model
	}{
		{
			name: "with none exported",
			val:  new(testModel),
			wantModel: &model{
				tableName: "`test_model`",
				columnMap: map[string]field{
					"id": {
						columnName: "id",
						fieldName:  "Id",
						typ:        reflect.TypeOf(uint64(0)),
					},
					"first_name": {
						columnName: "first_name",
						fieldName:  "FirstName",
						typ:        reflect.TypeOf(""),
					},
				},
				fieldMap: map[string]field{
					"Id": {
						columnName: "id",
						fieldName:  "Id",
						typ:        reflect.TypeOf(uint64(0)),
					},
					"FirstName": {
						columnName: "first_name",
						fieldName:  "FirstName",
						typ:        reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name:    "string",
			val:     "",
			wantErr: errors.New("数据类型不是指向结构体的指针: string"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRegistry()
			m, err := r.parseModel(reflect.TypeOf(tc.val))
			require.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

type testModelWithTableName struct {
	Id          uint64
	Name        string
	nonExported string
}

func (tm testModelWithTableName) TableName() string {
	return "`test_model`"
}

func modelWithTableName(tableName string) modelOpt {
	return func(model *model) error {
		model.tableName = tableName
		return nil
	}
}

func Test_tableName(t *testing.T) {
	type testModel struct {
	}
	testCases := []struct {
		name          string
		val           any
		opt           modelOpt
		wantErr       error
		wantTableName string
	}{
		{
			name:          "default table name",
			val:           new(testModel),
			wantTableName: "`test_model`",
		},
		{
			name:          "custom table name",
			val:           new(testModelWithTableName),
			wantTableName: "`test_model`",
		},
		{
			name:          "option table name",
			val:           new(testModelWithTableName),
			opt:           modelWithTableName("`test_model_t`"),
			wantTableName: "`test_model_t`",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRegistry()
			typ := reflect.TypeOf(tc.val)
			r.register(typ, tc.opt)
			m, err := r.get(typ)
			require.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantTableName, m.tableName)
		})
	}
}

func modelWithColumnName(fieldName, colName string) modelOpt {
	return func(model *model) error {
		if _, ok := model.fieldMap[fieldName]; !ok {
			return errors.New("字段名不存在: " + fieldName)
		}
		model.fieldMap[fieldName] = field{
			columnName: colName,
			fieldName:  fieldName,
		}
		return nil
	}
}

func Test_field(t *testing.T) {
	type testModel struct {
		Id        uint64
		FirstName string
	}
	type testModelWithTag struct {
		Id        uint64 `orm:"column=id"`
		FirstName string `orm:"column=first_name_c"`
	}
	testCases := []struct {
		name        string
		val         any
		opt         modelOpt
		fieldName   string
		wantErr     error
		wantColName string
	}{
		{
			name:        "default column name",
			val:         new(testModel),
			fieldName:   "FirstName",
			wantColName: "first_name",
		},
		{
			name:        "option column name",
			val:         new(testModel),
			opt:         modelWithColumnName("FirstName", "first_name_c"),
			fieldName:   "FirstName",
			wantColName: "first_name_c",
		},
		{
			name:    "option column name",
			val:     new(testModel),
			opt:     modelWithColumnName("FullName", "full_name_c"),
			wantErr: errors.New("字段名不存在: FullName"),
		},
		{
			name:        "with column tag",
			val:         new(testModelWithTag),
			fieldName:   "FirstName",
			wantColName: "first_name_c",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRegistry()
			typ := reflect.TypeOf(tc.val)
			m, err := r.register(typ, tc.opt)
			require.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantColName, m.fieldMap[tc.fieldName].columnName)
		})
	}
}
