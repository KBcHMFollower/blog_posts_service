package postgresrepository

import (
	"database/sql"
	"fmt"

	"github.com/KBcHMFollower/test_plate_user_service/cmd/migrator"
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

func (r *PostgresRepository) Migrate(pathToMigrates string) error {
	migrator, err := migrator.New(r.db)
	if err != nil {
		return fmt.Errorf("can`t create migrator : %v", err)
	}

	err = migrator.Migrate(pathToMigrates, "postgres")
	if err != nil {
		return fmt.Errorf("can`t migrate : %v", err)
	}

	return nil
}
