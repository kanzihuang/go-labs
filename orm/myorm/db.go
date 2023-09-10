package myorm

type DB struct {
	registry *registry
}

func NewDB() *DB {
	db := &DB{
		registry: newRegistry(),
	}
	return db
}
