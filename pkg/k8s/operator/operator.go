package operator

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/flack/chaos-engineering-as-a-platform/pkg/config"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/k8s"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/monitoring"
)

// ChaosOperator represents the chaos operator that manages chaos experiments
type ChaosOperator struct {
	client       *k8s.Client
	metrics      *monitoring.Metrics
	stopCh       chan struct{}
	wg           sync.WaitGroup
	experiments  map[string]ExperimentController
	experimentMu sync.RWMutex
}

// NewChaosOperator creates a new chaos operator
func NewChaosOperator(cfg *config.Config, metrics *monitoring.Metrics) (*ChaosOperator, error) {
	var client *k8s.Client
	var err error
	
	if cfg.MockKubernetes {
		log.Println("Using mock Kubernetes client for development")
		client = k8s.NewMockClient()
	} else {
		client, err = k8s.NewClient(cfg.KubeConfigPath, cfg.Namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
		}
	}

	return &ChaosOperator{
		client:      client,
		metrics:     metrics,
		stopCh:      make(chan struct{}),
		experiments: make(map[string]ExperimentController),
	}, nil
}

// Start starts the chaos operator
func (o *ChaosOperator) Start() error {
	log.Println("Starting chaos operator...")

	// Start the controller loop
	o.wg.Add(1)
	go o.controllerLoop()

	return nil
}

// Stop stops the chaos operator
func (o *ChaosOperator) Stop() error {
	log.Println("Stopping chaos operator...")

	// Signal all goroutines to stop
	close(o.stopCh)

	// Wait for all goroutines to finish
	o.wg.Wait()

	// Stop all running experiments
	o.experimentMu.Lock()
	for id, exp := range o.experiments {
		log.Printf("Stopping experiment %s", id)
		if err := exp.Stop(); err != nil {
			log.Printf("Error stopping experiment %s: %v", id, err)
		}
	}
	o.experimentMu.Unlock()

	return nil
}

// controllerLoop is the main loop of the operator
func (o *ChaosOperator) controllerLoop() {
	defer o.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-o.stopCh:
			log.Println("Controller loop stopping...")
			return
		case <-ticker.C:
			// Check for new experiments to run
			o.reconcileExperiments()
		}
	}
}

// reconcileExperiments checks for new experiments and updates existing ones
func (o *ChaosOperator) reconcileExperiments() {
	// This would typically involve checking a Custom Resource Definition (CRD)
	// for chaos experiments, but for simplicity, we'll just log
	log.Println("Reconciling experiments...")

	// In a real implementation, this would:
	// 1. List all ChaosExperiment CRs
	// 2. Start controllers for new experiments
	// 3. Update controllers for changed experiments
	// 4. Stop controllers for deleted experiments
}

// RunExperiment runs a new chaos experiment
func (o *ChaosOperator) RunExperiment(config *ExperimentConfig) error {
	o.experimentMu.Lock()
	defer o.experimentMu.Unlock()

	// Validate config
	if config == nil {
		return fmt.Errorf("experiment config cannot be nil")
	}

	// Check if experiment is already running
	if _, exists := o.experiments[config.ID]; exists {
		return fmt.Errorf("experiment %s is already running", config.ID)
	}

	// Determine if this is an external target
	isExternal := false
	if config.Params != nil && config.Params["target_type"] == "external" {
		isExternal = true
	}

	var controller ExperimentController

	if isExternal {
		// Create a new external experiment controller
		extController, err := NewExternalExperimentController(config.ID, config.Target, config.Params, config.Duration, o.metrics)
		if err != nil {
			return fmt.Errorf("failed to create external experiment controller: %w", err)
		}
		controller = extController
	} else {
		// Create a new Kubernetes experiment controller
		k8sController, err := NewExperimentController(config, o.client, o.metrics)
		if err != nil {
			return fmt.Errorf("failed to create experiment controller: %w", err)
		}
		controller = k8sController
	}

	// Start the experiment
	if err := controller.Start(); err != nil {
		return fmt.Errorf("failed to start experiment: %w", err)
	}

	// Store the controller
	o.experiments[config.ID] = controller

	// Update metrics
	o.metrics.ExperimentsExecuted.Inc()
	o.metrics.ActiveExperiments.Inc()

	return nil
}

// StopExperiment stops a running chaos experiment
func (o *ChaosOperator) StopExperiment(experimentID string) error {
	o.experimentMu.Lock()
	defer o.experimentMu.Unlock()

	// Check if experiment exists
	controller, exists := o.experiments[experimentID]
	if !exists {
		return fmt.Errorf("experiment %s not found", experimentID)
	}

	// Stop the experiment
	if err := controller.Stop(); err != nil {
		return fmt.Errorf("failed to stop experiment: %w", err)
	}

	// Remove the controller
	delete(o.experiments, experimentID)

	// Update metrics
	o.metrics.ActiveExperiments.Dec()

	return nil
}

// GetExperimentStatus gets the status of a running experiment
func (o *ChaosOperator) GetExperimentStatus(experimentID string) (string, error) {
	o.experimentMu.RLock()
	defer o.experimentMu.RUnlock()

	// Check if experiment exists
	controller, exists := o.experiments[experimentID]
	if !exists {
		return "", fmt.Errorf("experiment %s not found", experimentID)
	}

	return controller.GetStatus(), nil
}
