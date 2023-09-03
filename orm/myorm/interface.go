package myorm

type Query struct {
	SQL    string
	Params []any
}

type Expression interface {
	Build() (string, error)
}

type QueryBuilder interface {
	Build() (*Query, error)
	From(tableName string) QueryBuilder
	Where(expr Expression, args []any) QueryBuilder
}
