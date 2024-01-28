package builtin

type User struct {
}

type Order struct {
}

type Cache[T any] interface {
	Set(t T)
}
