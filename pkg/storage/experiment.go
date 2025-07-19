package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// ExperimentType defines the type of chaos experiment
type ExperimentType string

const (
	PodFailure     ExperimentType = "pod-failure"
	NetworkDelay   ExperimentType = "network-delay"
	CPUStress      ExperimentType = "cpu-stress"
	MemoryStress   ExperimentType = "memory-stress"
	DiskFailure    ExperimentType = "disk-failure"
	ServiceFailure ExperimentType = "service-failure"
	ExternalTarget ExperimentType = "external-target"
)

// ExperimentStatus defines the status of an experiment
type ExperimentStatus string

const (
	StatusPending   ExperimentStatus = "pending"
	StatusRunning   ExperimentStatus = "running"
	StatusCompleted ExperimentStatus = "completed"
	StatusFailed    ExperimentStatus = "failed"
	StatusCancelled ExperimentStatus = "cancelled"
)

// Experiment represents a chaos experiment
type Experiment struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        ExperimentType   `json:"type"`
	Status      ExperimentStatus `json:"status"`
	Target      string           `json:"target"`
	Parameters  string           `json:"parameters"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Duration    int              `json:"duration"` // Duration in seconds
}

// CreateExperiment creates a new experiment in the database
func (d *Database) CreateExperiment(experiment *Experiment) error {
	query := `
		INSERT INTO experiments (id, name, description, type, status, target, parameters, created_at, updated_at, duration)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	_, err := d.db.Exec(
		query,
		experiment.ID,
		experiment.Name,
		experiment.Description,
		experiment.Type,
		experiment.Status,
		experiment.Target,
		experiment.Parameters,
		experiment.CreatedAt,
		experiment.UpdatedAt,
		experiment.Duration,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create experiment: %w", err)
	}
	
	return nil
}

// GetExperiment retrieves an experiment by ID
func (d *Database) GetExperiment(id string) (*Experiment, error) {
	query := `
		SELECT id, name, description, type, status, target, parameters, created_at, updated_at, duration
		FROM experiments
		WHERE id = $1
	`
	
	var experiment Experiment
	err := d.db.QueryRow(query, id).Scan(
		&experiment.ID,
		&experiment.Name,
		&experiment.Description,
		&experiment.Type,
		&experiment.Status,
		&experiment.Target,
		&experiment.Parameters,
		&experiment.CreatedAt,
		&experiment.UpdatedAt,
		&experiment.Duration,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("experiment not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get experiment: %w", err)
	}
	
	return &experiment, nil
}

// ListExperiments retrieves all experiments
func (d *Database) ListExperiments() ([]*Experiment, error) {
	query := `
		SELECT id, name, description, type, status, target, parameters, created_at, updated_at, duration
		FROM experiments
		ORDER BY created_at DESC
	`
	
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list experiments: %w", err)
	}
	defer rows.Close()
	
	var experiments []*Experiment
	for rows.Next() {
		var experiment Experiment
		err := rows.Scan(
			&experiment.ID,
			&experiment.Name,
			&experiment.Description,
			&experiment.Type,
			&experiment.Status,
			&experiment.Target,
			&experiment.Parameters,
			&experiment.CreatedAt,
			&experiment.UpdatedAt,
			&experiment.Duration,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan experiment: %w", err)
		}
		experiments = append(experiments, &experiment)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating experiments: %w", err)
	}
	
	return experiments, nil
}

// UpdateExperimentStatus updates the status of an experiment
func (d *Database) UpdateExperimentStatus(id string, status ExperimentStatus) error {
	query := `
		UPDATE experiments
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	
	_, err := d.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update experiment status: %w", err)
	}
	
	return nil
}

// DeleteExperiment deletes an experiment by ID
func (d *Database) DeleteExperiment(id string) error {
	query := `
		DELETE FROM experiments
		WHERE id = $1
	`
	
	_, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete experiment: %w", err)
	}
	
	return nil
}