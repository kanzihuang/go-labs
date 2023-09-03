package myorm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModel struct {
	Id   int64
	Name string
}

func TestSelector_Build(t *testing.T) {
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:      "base",
			q:         Selector[TestModel]{},
			wantQuery: &Query{SQL: "SELECT * FROM `TestModel`;"},
		},
		{
			name:      "with FROM",
			q:         Selector[TestModel]{}.From("test_model"),
			wantQuery: &Query{SQL: "SELECT * FROM test_model;"},
		},
		{
			name:      "empty FROM",
			q:         Selector[TestModel]{}.From(""),
			wantQuery: &Query{SQL: "SELECT * FROM `TestModel`;"},
		},
		{
			name:      "with db",
			q:         Selector[TestModel]{}.From("test_db.test_model"),
			wantQuery: &Query{SQL: "SELECT * FROM test_db.test_model;"},
		},
		{
			name: "with where id",
			q: Selector[TestModel]{}.Where(Predicate{
				left:    C("id"),
				operate: OperateEqual,
				right:   Value{},
			}, []any{1}),
			wantQuery: &Query{
				SQL:    "SELECT * FROM `TestModel` WHERE `id` = ?;",
				Params: []any{1},
			},
		},
		{
			name: "with where id and name",
			q: Selector[TestModel]{}.Where(
				Predicate{
					left: Predicate{
						C("id"),
						OperateEqual,
						Value{},
					},
					operate: OperateAnd,
					right: Predicate{
						operate: OperateNot,
						right: Predicate{
							left: Predicate{
								C("name"),
								OperateEqual,
								Value{},
							},
							operate: OperateOr,
							right: Predicate{
								left:    C("name"),
								operate: OperateEqual,
								right:   Value{},
							},
						},
					},
				},
				[]any{1, "a", "b"}),
			wantQuery: &Query{
				SQL:    "SELECT * FROM `TestModel` WHERE `id` = ? AND NOT (`name` = ? OR `name` = ?);",
				Params: []any{1, "a", "b"},
			},
		},
		{
			name: "with where not name",
			q: Selector[TestModel]{}.Where(
				Predicate{
					operate: OperateNot,
					right: Predicate{
						Predicate{
							C("name"),
							OperateEqual,
							Value{},
						},
						OperateOr,
						Predicate{
							C("name"),
							OperateEqual,
							Value{},
						},
					},
				},
				[]any{1, "a", "b"}),
			wantQuery: &Query{
				SQL:    "SELECT * FROM `TestModel` WHERE NOT (`name` = ? OR `name` = ?);",
				Params: []any{1, "a", "b"},
			},
		},
	}
	for _, tc := range testCases {
		query, err := tc.q.Build()
		assert.Equal(t, tc.wantErr, err)
		if err != nil {
			return
		}
		assert.Equal(t, tc.wantQuery, query)
	}
}
