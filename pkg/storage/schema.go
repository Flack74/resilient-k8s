package storage

import (
	"fmt"
)

// InitSchema initializes the database schema
func (d *Database) InitSchema() error {
	// Create experiments table
	experimentsTable := `
		CREATE TABLE IF NOT EXISTS experiments (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			type VARCHAR(50) NOT NULL,
			status VARCHAR(50) NOT NULL,
			target VARCHAR(255) NOT NULL,
			parameters JSONB,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			duration INTEGER NOT NULL
		)
	`
	
	// Create targets table
	targetsTable := `
		CREATE TABLE IF NOT EXISTS targets (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			type VARCHAR(50) NOT NULL,
			namespace VARCHAR(255) NOT NULL,
			selector TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`
	
	// Create experiment results table
	resultsTable := `
		CREATE TABLE IF NOT EXISTS experiment_results (
			id VARCHAR(36) PRIMARY KEY,
			experiment_id VARCHAR(36) NOT NULL,
			status VARCHAR(50) NOT NULL,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP,
			metrics JSONB,
			logs TEXT,
			FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE
		)
	`
	
	// Execute the schema creation
	if _, err := d.db.Exec(experimentsTable); err != nil {
		return fmt.Errorf("failed to create experiments table: %w", err)
	}
	
	if _, err := d.db.Exec(targetsTable); err != nil {
		return fmt.Errorf("failed to create targets table: %w", err)
	}
	
	if _, err := d.db.Exec(resultsTable); err != nil {
		return fmt.Errorf("failed to create experiment_results table: %w", err)
	}
	
	return nil
}