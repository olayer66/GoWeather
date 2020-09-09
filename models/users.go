package models

import (
	"database/sql"
	"errors"
)

type User struct {
	Id       int
	Name     string
	Surname  string
	Usuario  string
	Password string
	ApiKey   string
}

func GetGetUserById(db *sql.DB, id int) (string, error) {
	var s sql.NullString
	sql := "SELECT * FROM users WHERE Id='" + string(id) + "'"
	err := db.QueryRow(sql).Scan(&s)
	if err != nil {
		return "", err
	}
	if s.Valid {
		return s.String, nil
	} else {
		return "", errors.New("No rows returned")
	}
}
