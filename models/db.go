package models

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Datastore interface {
	GetUserById(db *DB, Id int) ([]*User, error)
}

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
