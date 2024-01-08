package myorm

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	private   int64
}

func TestSelector_Build(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		q         *Selector[TestModel]
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
								C("FirstName"),
								OperateEqual,
								Value{},
							},
							operate: OperateOr,
							right: Predicate{
								left:    C("FirstName"),
								operate: OperateEqual,
								right:   Value{},
							},
						},
					},
				},
				[]any{1, "a", "b"}),
			wantQuery: &Query{
				SQL:    "SELECT * FROM `test_model` WHERE `id` = ? AND NOT (`first_name` = ? OR `first_name` = ?);",
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
							C("FirstName"),
							OperateEqual,
							Value{},
						},
						OperateOr,
						Predicate{
							C("FirstName"),
							OperateEqual,
							Value{},
						},
					},
				},
				[]any{1, "a", "b"}),
			wantQuery: &Query{
				SQL:    "SELECT * FROM `test_model` WHERE NOT (`first_name` = ? OR `first_name` = ?);",
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

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		selector *Selector[TestModel]
		mockRows *sqlmock.Rows
		mockErr  error
		wantData *TestModel
		wantErr  error
	}{
		{
			name:     "select",
			selector: NewSelector[TestModel](db),
			mockRows: sqlmock.NewRows([]string{"id", "first_name"}).
				AddRow(1, "Mike"),
			wantData: &TestModel{
				Id:        1,
				FirstName: "Mike",
			},
		},
		{
			name: "select where invalid",
			selector: NewSelector[TestModel](db).Where(Predicate{
				left:    C("invalid"),
				operate: OperateEqual,
				right:   Value{},
			}, []any{1}),
			mockErr: errors.New("invalid column: invalid"),
			wantErr: errors.New("invalid column: invalid"),
		},
	}
	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			eq := mock.ExpectQuery("SELECT .*")
			if tc.mockErr != nil {
				eq.WillReturnError(tc.mockErr)
			} else {
				eq.WillReturnRows(tc.mockRows)
			}
			data, err := tc.selector.Get(ctx)
			require.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantData, data)
		})
	}
}
