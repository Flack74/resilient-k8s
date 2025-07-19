package operator

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/flack/chaos-engineering-as-a-platform/pkg/monitoring"
)

// ExternalExperimentController controls chaos experiments on external targets
type ExternalExperimentController struct {
	experimentID string
	targetURL    string
	params       map[string]string
	duration     int
	client       *http.Client
	metrics      *monitoring.Metrics
	stopCh       chan struct{}
	status       string
}

// NewExternalExperimentController creates a new external experiment controller
func NewExternalExperimentController(experimentID, targetURL string, params map[string]string, duration int, metrics *monitoring.Metrics) (*ExternalExperimentController, error) {
	return &ExternalExperimentController{
		experimentID: experimentID,
		targetURL:    targetURL,
		params:       params,
		duration:     duration,
		client:       &http.Client{Timeout: 10 * time.Second},
		metrics:      metrics,
		stopCh:       make(chan struct{}),
		status:       "pending",
	}, nil
}

// Start starts the external experiment
func (c *ExternalExperimentController) Start() error {
	log.Printf("Starting external experiment %s on target %s", c.experimentID, c.targetURL)
	
	// Update status
	c.status = "running"
	
	// Start the experiment in a goroutine
	go c.runExperiment()
	
	return nil
}

// Stop stops the external experiment
func (c *ExternalExperimentController) Stop() error {
	log.Printf("Stopping external experiment %s", c.experimentID)
	
	// Signal the goroutine to stop
	close(c.stopCh)
	
	// Update status
	c.status = "stopped"
	
	return nil
}

// GetStatus gets the status of the external experiment
func (c *ExternalExperimentController) GetStatus() string {
	return c.status
}

// runExperiment runs the external experiment
func (c *ExternalExperimentController) runExperiment() {
	// Create a timer for the experiment duration
	timer := time.NewTimer(time.Duration(c.duration) * time.Second)
	defer timer.Stop()
	
	// Execute the experiment
	if err := c.executeExperiment(); err != nil {
		log.Printf("Error executing external experiment %s: %v", c.experimentID, err)
		c.status = "failed"
		return
	}
	
	// Wait for the experiment to complete or be stopped
	select {
	case <-timer.C:
		// Experiment completed successfully
		log.Printf("External experiment %s completed", c.experimentID)
		c.status = "completed"
	case <-c.stopCh:
		// Experiment was stopped
		log.Printf("External experiment %s was stopped", c.experimentID)
		c.status = "stopped"
	}
	
	// Clean up the experiment
	if err := c.cleanupExperiment(); err != nil {
		log.Printf("Error cleaning up external experiment %s: %v", c.experimentID, err)
	}
}

// executeExperiment executes the external experiment
func (c *ExternalExperimentController) executeExperiment() error {
	// Create the request URL with parameters
	url := c.targetURL
	if c.params["endpoint"] != "" {
		url += c.params["endpoint"]
	}
	
	// Create the request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Chaos-Experiment-ID", c.experimentID)
	req.Header.Set("X-Chaos-Experiment-Type", c.params["type"])
	req.Header.Set("X-Chaos-Experiment-Duration", fmt.Sprintf("%d", c.duration))
	
	// Add authorization if provided
	if c.params["auth_token"] != "" {
		req.Header.Set("Authorization", "Bearer "+c.params["auth_token"])
	}
	
	// Execute the request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check the response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	return nil
}

// cleanupExperiment cleans up the external experiment
func (c *ExternalExperimentController) cleanupExperiment() error {
	// Create the request URL with parameters
	url := c.targetURL
	if c.params["cleanup_endpoint"] != "" {
		url += c.params["cleanup_endpoint"]
	} else {
		url += "/cleanup"
	}
	
	// Create the request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create cleanup request: %w", err)
	}
	
	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Chaos-Experiment-ID", c.experimentID)
	
	// Add authorization if provided
	if c.params["auth_token"] != "" {
		req.Header.Set("Authorization", "Bearer "+c.params["auth_token"])
	}
	
	// Execute the request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute cleanup request: %w", err)
	}
	defer resp.Body.Close()
	
	return nil
}