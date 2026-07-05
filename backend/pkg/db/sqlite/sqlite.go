package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3" // Driver initialization
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Database struct {
	Db *sql.DB
}

// Connect handles defensive configuration optimization parameters for SQLite
func Connect(dbPath string) (*Database, error) {
	// Defensive tuning parameters:
	// _foreign_keys=on - Enforces relational integrity constraints at the engine level
	// _journal_mode=WAL - Write-Ahead Logging allows concurrent reads without locking writes
	// _busy_timeout=5000 - Wait up to 5000ms if DB is locked before returning an error
	dsn := fmt.Sprintf("%s?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000", dbPath)
	
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Verify physical communication link to the storage medium
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{Db: db}, nil
}

// RunMigrations compiles and executes migrations atomically using io/fs virtual directory
func (d *Database) RunMigrations() error {
	sourceDriver, err := iofs.New(migrationFiles, "migrations/sqlite")
	if err != nil {
		return fmt.Errorf("failed to create migration iofs driver: %w", err)
	}

	dbDriver, err := sqlite3.WithInstance(d.Db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to instantiate sqlite3 migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", dbDriver)
	if err != nil {
		return fmt.Errorf("failed to initialize migration controller: %w", err)
	}

	// Apply migrations up to the current newest version
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration run failed: %w", err)
	}

	log.Println("[DATABEDROCK] Database schemas migrated successfully.")
	return nil
}