package sqlite

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/sqlite/*.sql
var migrationFiles embed.FS

type Database struct {
	DB *sql.DB
}

// Connect establishes a connection to the SQLite database and returns a Database instance.
func Connect(dbPath string) (*Database, error) {
	dsn := fmt.Sprintf("%s?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000", dbPath)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	// Verify the connection by pinging the database and checking for any errors.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	return &Database{DB: db}, nil
}

// RunMigrations executes all pending database migrations against the connected SQLite database.
func (d *Database) RunMigrations() error {
	sourceDriver, err := iofs.New(migrationFiles, "migrations/sqlite")
	if err != nil {
		return fmt.Errorf("failed to create migration source driver: %w", err)
	}

	dbDriver, err := sqlite3.WithInstance(d.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create SQLite database driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", dbDriver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

