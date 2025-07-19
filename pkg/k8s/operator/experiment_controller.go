package operator

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/flack/chaos-engineering-as-a-platform/pkg/k8s"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/monitoring"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/storage"
)

// ExperimentConfig holds configuration for an experiment
type ExperimentConfig struct {
	ID             string
	Type           string
	Target         string
	Params         map[string]string
	Duration       int
}

// K8sExperimentController controls a single chaos experiment in Kubernetes
type K8sExperimentController struct {
	id             string
	experimentType string
	target         string
	params         map[string]string
	duration       int
	client         *k8s.Client
	metrics        *monitoring.Metrics
	status         storage.ExperimentStatus
	statusMu       sync.RWMutex
	stopCh         chan struct{}
	doneCh         chan struct{}
}

// NewExperimentController creates a new experiment controller
func NewExperimentController(config *ExperimentConfig, client *k8s.Client, metrics *monitoring.Metrics) (*K8sExperimentController, error) {
	if config == nil {
		return nil, fmt.Errorf("experiment config cannot be nil")
	}
	
	return &K8sExperimentController{
		id:             config.ID,
		experimentType: config.Type,
		target:         config.Target,
		params:         config.Params,
		duration:       config.Duration,
		client:         client,
		metrics:        metrics,
		status:         storage.StatusPending,
		stopCh:         make(chan struct{}),
		doneCh:         make(chan struct{}),
	}, nil
}

// Start starts the experiment
func (c *K8sExperimentController) Start() error {
	c.setStatus(storage.StatusRunning)

	go func() {
		defer close(c.doneCh)

		// Execute the experiment based on its type
		var err error
		switch c.experimentType {
		case "pod-failure":
			err = c.executePodFailure()
		case "network-delay":
			err = c.executeNetworkDelay()
		case "cpu-stress":
			err = c.executeCPUStress()
		case "memory-stress":
			err = c.executeMemoryStress()
		default:
			err = fmt.Errorf("unsupported experiment type: %s", c.experimentType)
		}

		if err != nil {
			log.Printf("Experiment %s failed: %v", c.id, err)
			c.setStatus(storage.StatusFailed)
			c.metrics.ExperimentsFailed.Inc()
		} else {
			c.setStatus(storage.StatusCompleted)
			c.metrics.ExperimentsSucceeded.Inc()
		}
	}()

	return nil
}

// Stop stops the experiment
func (c *K8sExperimentController) Stop() error {
	// Signal the experiment to stop
	close(c.stopCh)

	// Wait for the experiment to finish
	<-c.doneCh

	// Update status if it was running
	c.statusMu.Lock()
	if c.status == storage.StatusRunning {
		c.status = storage.StatusCancelled
	}
	c.statusMu.Unlock()

	return nil
}

// GetStatus gets the current status of the experiment
func (c *K8sExperimentController) GetStatus() string {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	return string(c.status)
}

// setStatus sets the status of the experiment
func (c *K8sExperimentController) setStatus(status storage.ExperimentStatus) {
	c.statusMu.Lock()
	c.status = status
	c.statusMu.Unlock()
}

// executeExperiment is a generic method to execute experiments
func (c *K8sExperimentController) executeExperiment(experimentType string) error {
	log.Printf("Executing %s experiment on target %s", experimentType, c.target)

	// Simulate the experiment duration
	select {
	case <-time.After(time.Duration(c.duration) * time.Second):
		log.Printf("%s experiment completed for target %s", experimentType, c.target)
	case <-c.stopCh:
		log.Printf("%s experiment cancelled for target %s", experimentType, c.target)
		return nil
	}

	return nil
}

// executePodFailure executes a pod failure experiment
func (c *K8sExperimentController) executePodFailure() error {
	// In a real implementation, this would use the Kubernetes API to find and delete pods
	// based on the target selector
	return c.executeExperiment("pod failure")
}

// executeNetworkDelay executes a network delay experiment
func (c *K8sExperimentController) executeNetworkDelay() error {
	// In a real implementation, this would use network policies or a CNI plugin
	// to introduce network delays
	return c.executeExperiment("network delay")
}

// executeCPUStress executes a CPU stress experiment
func (c *K8sExperimentController) executeCPUStress() error {
	// In a real implementation, this would deploy a stress testing pod
	// to consume CPU resources
	return c.executeExperiment("CPU stress")
}

// executeMemoryStress executes a memory stress experiment
func (c *K8sExperimentController) executeMemoryStress() error {
	// In a real implementation, this would deploy a stress testing pod
	// to consume memory resources
	return c.executeExperiment("memory stress")
}
