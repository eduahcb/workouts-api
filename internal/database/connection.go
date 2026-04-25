package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Open(logger *log.Logger) (*pgxpool.Pool, error) {

	url := "postgres://postgres:postgres@localhost:5432/postgres"
	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("db: parse config %w", err)
	}

	logger.Println("Connected to Database")

	return conn, nil
}

func MigrateFS(pool *pgxpool.Pool, migrationsFS fs.FS, dir string, logger *log.Logger) error {
	config := pool.Config().ConnConfig

	db := stdlib.OpenDB(*config)
	defer db.Close()

	goose.SetBaseFS(migrationsFS)
	defer goose.SetBaseFS(nil)

	return Migrate(db, dir, logger)
}

func Migrate(db *sql.DB, dir string, logger *log.Logger) error {
	err := goose.SetDialect("postgres")

	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)

	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	logger.Println("Migrations applied successfully")

	return nil
}
