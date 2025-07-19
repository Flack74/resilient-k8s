package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// TargetType defines the type of target for chaos experiments
type TargetType string

const (
	TargetPod       TargetType = "pod"
	TargetDeployment TargetType = "deployment"
	TargetService    TargetType = "service"
	TargetNode       TargetType = "node"
)

// Target represents a target for chaos experiments
type Target struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        TargetType `json:"type"`
	Namespace   string     `json:"namespace"`
	Selector    string     `json:"selector"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateTarget creates a new target in the database
func (d *Database) CreateTarget(target *Target) error {
	query := `
		INSERT INTO targets (id, name, description, type, namespace, selector, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := d.db.Exec(
		query,
		target.ID,
		target.Name,
		target.Description,
		target.Type,
		target.Namespace,
		target.Selector,
		target.CreatedAt,
		target.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create target: %w", err)
	}
	
	return nil
}

// GetTarget retrieves a target by ID
func (d *Database) GetTarget(id string) (*Target, error) {
	query := `
		SELECT id, name, description, type, namespace, selector, created_at, updated_at
		FROM targets
		WHERE id = $1
	`
	
	var target Target
	err := d.db.QueryRow(query, id).Scan(
		&target.ID,
		&target.Name,
		&target.Description,
		&target.Type,
		&target.Namespace,
		&target.Selector,
		&target.CreatedAt,
		&target.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("target not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get target: %w", err)
	}
	
	return &target, nil
}

// ListTargets retrieves all targets
func (d *Database) ListTargets() ([]*Target, error) {
	query := `
		SELECT id, name, description, type, namespace, selector, created_at, updated_at
		FROM targets
		ORDER BY created_at DESC
	`
	
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list targets: %w", err)
	}
	defer rows.Close()
	
	var targets []*Target
	for rows.Next() {
		var target Target
		err := rows.Scan(
			&target.ID,
			&target.Name,
			&target.Description,
			&target.Type,
			&target.Namespace,
			&target.Selector,
			&target.CreatedAt,
			&target.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan target: %w", err)
		}
		targets = append(targets, &target)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating targets: %w", err)
	}
	
	return targets, nil
}

// DeleteTarget deletes a target by ID
func (d *Database) DeleteTarget(id string) error {
	query := `
		DELETE FROM targets
		WHERE id = $1
	`
	
	_, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete target: %w", err)
	}
	
	return nil
}