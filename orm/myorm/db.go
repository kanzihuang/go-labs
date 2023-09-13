package myorm

import "database/sql"

type DB struct {
	db           *sql.DB
	registry     *registry
	valueCreator valueCreator
}

type DBOption func(db *DB) error

func UseRedirectValue(db *DB) {
	db.valueCreator = newReflectValue
}

func Open(driverName string, dataSourceName string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	myDB := &DB{
		db:           db,
		registry:     newRegistry(),
		valueCreator: newUnsafeValue,
	}
	for _, opt := range opts {
		if err := opt(myDB); err != nil {
			return nil, err
		}
	}
	return myDB, nil
}
