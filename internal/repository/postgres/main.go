package postgresrepository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func New(connecionString string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connecionString)
	if err != nil {
		return nil, fmt.Errorf("Error in process db connection : %v", err)
	}

	return &PostgresRepository{
		db: db,
	}, nil
}
