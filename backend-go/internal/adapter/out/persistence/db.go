package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string, migrationsPath string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s?_foreign_keys=on", dbPath)
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("error trying to open Sqlite connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error to connect to Sqlite: %w", err)
	}

	if err := runMigrations(db.DB, migrationsPath); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error trying to apply migrations: %w", err)
	}

	log.Println("SQLite database connected and migrations applied")
	return db, nil
}

func runMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("error trying to create Sqlite driver for migrations: %w", err)
	}

	sourceURL := fmt.Sprintf("file://%s", migrationsPath)
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"sqlite3",
		driver,
	)

	if err != nil {
		return fmt.Errorf("error trying to create migrations: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error trying to run migrations: %w", err)
	}

	return nil
}
