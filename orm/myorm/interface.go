package myorm

type TableName interface {
	TableName() string
}

type Query struct {
	SQL    string
	Params []any
}

type Expression interface {
	Build(mdl *model) (string, error)
}

type QueryBuilder interface {
	Build() (*Query, error)
	Where(expr Expression, args []any) QueryBuilder
}
