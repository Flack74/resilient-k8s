package experiments

import (
	"time"
)

// ExperimentResult represents the result of a chaos experiment
type ExperimentResult struct {
	ID                string    `json:"id"`
	ExperimentID      string    `json:"experiment_id"`
	ExperimentType    string    `json:"experiment_type"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	Duration          float64   `json:"duration"`
	Success           bool      `json:"success"`
	Error             string    `json:"error,omitempty"`
	AffectedResources []string  `json:"affected_resources,omitempty"`
	Metrics           map[string]float64 `json:"metrics,omitempty"`
}

// CalculateDuration calculates the duration of the experiment
func (r *ExperimentResult) CalculateDuration() {
	if !r.EndTime.IsZero() {
		r.Duration = r.EndTime.Sub(r.StartTime).Seconds()
	}
}