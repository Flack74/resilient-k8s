package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Database represents a database connection
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(connectionURL string) (*Database, error) {
	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// GetDB returns the underlying database connection
func (d *Database) GetDB() *sql.DB {
	return d.db
}