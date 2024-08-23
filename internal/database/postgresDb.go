package database

import (
	"database/sql"
	"fmt"
)

type DBDriver struct {
	*sql.DB
}

func NewPostgresDriver(connectionString string) (*DBDriver, *sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, nil, fmt.Errorf("error in process db connection : %v", err)
	}
	return &DBDriver{db}, db, nil
}
