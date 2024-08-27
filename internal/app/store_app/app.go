package store_app

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/KBcHMFollower/blog_posts_service/internal/config"
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type PostgresStore struct {
	Store         database.DBWrapper
	migrationPath string
	db            *sql.DB
} //TODO: ПЕРЕПИСАТЬ МИГРАТОР, ЭТО НЕ НОРМ

type StoreApp struct {
	PostgresStore *PostgresStore
}

func New(postgresConnectionInfo config.Storage) (*StoreApp, error) {
	db, err := sql.Open("postgres", postgresConnectionInfo.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error in process db connection : %w", err)
	}
	sqlxDb := sqlx.NewDb(db, "postgres")

	return &StoreApp{
		PostgresStore: &PostgresStore{
			Store:         &database.DBDriver{sqlxDb},
			migrationPath: postgresConnectionInfo.MigrationPath,
			db:            db,
		},
	}, nil
}

func (app *StoreApp) Run() error {
	if err := database.ForceMigrate(
		app.PostgresStore.db,
		app.PostgresStore.migrationPath,
		database.MigrateUp,
		database.Postgres,
	); err != nil {
		log.Fatalf("error in process db connection : %v", err)
		return err
	}

	return nil
}

func (app *StoreApp) Stop() error {
	var resErr error = nil

	if err := app.PostgresStore.Store.Stop(); err != nil {
		resErr = errors.Join(resErr, fmt.Errorf("error in stop postgres store : %w", err))
	}

	return resErr
}
