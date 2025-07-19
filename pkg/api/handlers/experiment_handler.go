package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/flack/chaos-engineering-as-a-platform/pkg/k8s/operator"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/monitoring"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/storage"
)

// ExperimentHandler handles experiment-related API requests
type ExperimentHandler struct {
	db      *storage.Database
	metrics *monitoring.Metrics
	operator *operator.ChaosOperator
}

// NewExperimentHandler creates a new experiment handler
func NewExperimentHandler(db *storage.Database, metrics *monitoring.Metrics) *ExperimentHandler {
	return &ExperimentHandler{
		db:      db,
		metrics: metrics,
	}
}

// SetOperator sets the chaos operator
func (h *ExperimentHandler) SetOperator(operator *operator.ChaosOperator) {
	h.operator = operator
}

// CreateExperimentRequest represents a request to create a new experiment
type CreateExperimentRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Type        string            `json:"type" binding:"required"`
	Target      string            `json:"target" binding:"required"`
	Parameters  map[string]string `json:"parameters"`
	Duration    int               `json:"duration" binding:"required"`
}

// CreateExperiment handles the creation of a new experiment
func (h *ExperimentHandler) CreateExperiment(c *gin.Context) {
	var req CreateExperimentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert parameters map to JSON string
	paramsJSON, err := json.Marshal(req.Parameters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal parameters"})
		return
	}

	// Create a new experiment
	now := time.Now()
	experiment := &storage.Experiment{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Type:        storage.ExperimentType(req.Type),
		Status:      storage.StatusPending,
		Target:      req.Target,
		Parameters:  string(paramsJSON),
		CreatedAt:   now,
		UpdatedAt:   now,
		Duration:    req.Duration,
	}

	// Save the experiment to the database
	if err := h.db.CreateExperiment(experiment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update metrics
	h.metrics.ExperimentsCreated.Inc()

	c.JSON(http.StatusCreated, experiment)
}

// ListExperiments handles listing all experiments
func (h *ExperimentHandler) ListExperiments(c *gin.Context) {
	experiments, err := h.db.ListExperiments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, experiments)
}

// GetExperiment handles retrieving a single experiment
func (h *ExperimentHandler) GetExperiment(c *gin.Context) {
	id := c.Param("id")
	experiment, err := h.db.GetExperiment(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, experiment)
}

// ExecuteExperiment handles executing an experiment
func (h *ExperimentHandler) ExecuteExperiment(c *gin.Context) {
	id := c.Param("id")
	
	// Get the experiment from the database
	experiment, err := h.db.GetExperiment(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Check if the operator is available
	if h.operator == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "chaos operator not available"})
		return
	}

	// Update the experiment status
	if err := h.db.UpdateExperimentStatus(id, storage.StatusRunning); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse parameters from JSON string
	var params map[string]string
	if err := json.Unmarshal([]byte(experiment.Parameters), &params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse parameters"})
		return
	}

	// Create experiment config
	config := &operator.ExperimentConfig{
		ID:       experiment.ID,
		Type:     string(experiment.Type),
		Target:   experiment.Target,
		Params:   params,
		Duration: experiment.Duration,
	}

	// Run the experiment
	if err := h.operator.RunExperiment(config); err != nil {
		// Revert the status if the experiment fails to start
		_ = h.db.UpdateExperimentStatus(id, storage.StatusFailed)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     experiment.ID,
		"status": "running",
	})
}

// DeleteExperiment handles deleting an experiment
func (h *ExperimentHandler) DeleteExperiment(c *gin.Context) {
	id := c.Param("id")
	
	// Check if the experiment exists
	_, err := h.db.GetExperiment(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Delete the experiment
	if err := h.db.DeleteExperiment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "experiment deleted"})
}