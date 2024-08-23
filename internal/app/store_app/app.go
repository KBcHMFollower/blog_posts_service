package store_app

import (
	"database/sql"
	"fmt"
	"github.com/KBcHMFollower/blog_posts_service/internal/config"
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
	_ "github.com/lib/pq"
	"log"
)

type PostgresStore struct {
	Store         database.DBWrapper
	migrationPath string
	db            *sql.DB
} //TODO: ПЕРЕПИСАТЬ МИГРАТОР, ЭТО НЕ НОРМ, хотя скорее переписать структуру бд

type StoreApp struct {
	PostgresStore *PostgresStore
}

func New(postgresConnectionInfo config.Storage) (*StoreApp, error) {
	db, err := sql.Open("postgres", postgresConnectionInfo.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error in process db connection : %w", err)
	}

	return &StoreApp{
		PostgresStore: &PostgresStore{
			Store:         &database.DBDriver{db},
			migrationPath: postgresConnectionInfo.MigrationPath,
			db:            db,
		},
	}, nil
}

func (app *StoreApp) Run() error {
	if err := database.ForceMigrate(app.PostgresStore.db, app.PostgresStore.migrationPath); err != nil {
		log.Fatalf("error in process db connection : %w", err)
		return err
	}

	return nil
}

func (app *StoreApp) Stop() error {
	if err := app.PostgresStore.Store.Close(); err != nil {
		return fmt.Errorf("error in stop postgres store_app : %w", err)
	}

	return nil
} //TODO: ДОЛЖНО СОБИРАТЬ СТЕК ОШИБОК, А НЕ ЗАВЕРШАТЬСЯ, КОГДА СЛОВИЛА ОДНУ
