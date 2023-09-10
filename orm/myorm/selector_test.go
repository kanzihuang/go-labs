package myorm

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestModel struct {
	Id   int64
	Name string
}

func TestSelector_Build(t *testing.T) {
	db := NewDB()
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:      "base",
			q:         NewSelector[TestModel](db),
			wantQuery: &Query{SQL: "SELECT * FROM `test_model`;"},
		},
		{
			name: "with where id",
			q: NewSelector[TestModel](db).Where(Predicate{
				left:    C("Id"),
				operate: OperateEqual,
				right:   Value{},
			}, []any{1}),
			wantQuery: &Query{
				SQL:    "SELECT * FROM `test_model` WHERE `id` = ?;",
				Params: []any{1},
			},
		},
		{
			name: "with where id and name",
			q: NewSelector[TestModel](db).Where(
				Predicate{
					left: Predicate{
						C("Id"),
						OperateEqual,
						Value{},
					},
					operate: OperateAnd,
					right: Predicate{
						operate: OperateNot,
						right: Predicate{
							left: Predicate{
								C("Name"),
								OperateEqual,
								Value{},
							},
							operate: OperateOr,
							right: Predicate{
								left:    C("Name"),
								operate: OperateEqual,
								right:   Value{},
							},
						},
					},
				},
				[]any{1, "a", "b"}),
			wantQuery: &Query{
				SQL:    "SELECT * FROM `test_model` WHERE `id` = ? AND NOT (`name` = ? OR `name` = ?);",
				Params: []any{1, "a", "b"},
			},
		},
		{
			name: "with where not name",
			q: NewSelector[TestModel](db).Where(
				Predicate{
					operate: OperateNot,
					right: Predicate{
						Predicate{
							C("Name"),
							OperateEqual,
							Value{},
						},
						OperateOr,
						Predicate{
							C("Name"),
							OperateEqual,
							Value{},
						},
					},
				},
				[]any{1, "a", "b"}),
			wantQuery: &Query{
				SQL:    "SELECT * FROM `test_model` WHERE NOT (`name` = ? OR `name` = ?);",
				Params: []any{1, "a", "b"},
			},
		},
		{
			name: "with where invalid column",
			q: NewSelector[TestModel](db).Where(Predicate{
				left:    C("invalid"),
				operate: OperateEqual,
				right:   Value{},
			}, []any{1}),
			wantErr: errors.New("invalid column: invalid"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			require.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
