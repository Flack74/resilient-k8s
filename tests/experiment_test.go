package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/flack/chaos-engineering-as-a-platform/pkg/storage"
	"github.com/google/uuid"
)

func TestExperiment(t *testing.T) {
	// Create parameters as JSON string
	params := map[string]string{"namespace": "default", "percentage": "50"}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal parameters: %v", err)
	}

	// Create a test experiment
	experiment := &storage.Experiment{
		ID:          uuid.New().String(),
		Name:        "Test Experiment",
		Description: "A test experiment",
		Type:        storage.PodFailure,
		Status:      storage.StatusPending,
		Target:      "app=test",
		Parameters:  string(paramsJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Duration:    60,
	}

	// Validate experiment fields
	if experiment.Name != "Test Experiment" {
		t.Errorf("Expected name to be 'Test Experiment', got '%s'", experiment.Name)
	}

	if experiment.Type != storage.PodFailure {
		t.Errorf("Expected type to be 'pod-failure', got '%s'", experiment.Type)
	}

	if experiment.Status != storage.StatusPending {
		t.Errorf("Expected status to be 'pending', got '%s'", experiment.Status)
	}

	if experiment.Duration != 60 {
		t.Errorf("Expected duration to be 60, got %d", experiment.Duration)
	}

	// Parse parameters from JSON string
	var parsedParams map[string]string
	if err := json.Unmarshal([]byte(experiment.Parameters), &parsedParams); err != nil {
		t.Fatalf("Failed to unmarshal parameters: %v", err)
	}

	// Test parameter access
	namespace, ok := parsedParams["namespace"]
	if !ok || namespace != "default" {
		t.Errorf("Expected namespace parameter to be 'default', got '%s'", namespace)
	}

	percentage, ok := parsedParams["percentage"]
	if !ok || percentage != "50" {
		t.Errorf("Expected percentage parameter to be '50', got '%s'", percentage)
	}
}