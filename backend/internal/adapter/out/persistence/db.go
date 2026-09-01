package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB opens a SQLite connection and applies migrations from a filesystem path.
// Used in development and tests where migrations live on disk.
func InitDB(dbPath string, migrationsPath string) (*sqlx.DB, error) {
	db, err := openDB(dbPath)
	if err != nil {
		return nil, err
	}

	if err := runMigrationsFromPath(db.DB, migrationsPath); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error trying to apply migrations: %w", err)
	}

	log.Println("SQLite database connected and migrations applied")
	return db, nil
}

// InitDBWithFS opens a SQLite connection and applies migrations from an embedded filesystem.
// Used in production where migrations are compiled into the binary via go:embed.
func InitDBWithFS(dbPath string, migrationsFS fs.FS) (*sqlx.DB, error) {
	db, err := openDB(dbPath)
	if err != nil {
		return nil, err
	}

	if err := runMigrationsFromFS(db.DB, migrationsFS); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error trying to apply migrations: %w", err)
	}

	log.Println("SQLite database connected and migrations applied")
	return db, nil
}

func openDB(dbPath string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s?_foreign_keys=on&_busy_timeout=5000", dbPath)
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("error trying to open Sqlite connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error to connect to Sqlite: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error enabling SQLite WAL mode: %w", err)
	}

	return db, nil
}

func runMigrationsFromPath(db *sql.DB, migrationsPath string) error {
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

func runMigrationsFromFS(db *sql.DB, migrationsFS fs.FS) error {
	dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("error trying to create Sqlite driver for migrations: %w", err)
	}

	sourceDriver, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("error trying to create embedded migrations source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", dbDriver)
	if err != nil {
		return fmt.Errorf("error trying to create migrations: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error trying to run migrations: %w", err)
	}

	return nil
}
