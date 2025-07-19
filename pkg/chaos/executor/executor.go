package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/flack/chaos-engineering-as-a-platform/pkg/chaos/experiments"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/k8s"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/monitoring"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/storage"
)

// ExperimentParams holds the common parameters for all experiments
type ExperimentParams struct {
	Namespace string
	Selector  string
	Value     int    // Generic value (percentage, delay, load, size)
	ValueUnit string // Unit for the value (%, ms, MB)
}

// Executor executes chaos experiments by running them against targets and tracking their results
type Executor struct {
	client  *k8s.Client
	db      *storage.Database
	metrics *monitoring.Metrics
}

// NewExecutor creates a new experiment executor with the provided Kubernetes client, database, and metrics
func NewExecutor(client *k8s.Client, db *storage.Database, metrics *monitoring.Metrics) *Executor {
	if client == nil {
		panic("k8s client cannot be nil")
	}
	if db == nil {
		panic("database cannot be nil")
	}
	if metrics == nil {
		panic("metrics cannot be nil")
	}

	return &Executor{
		client:  client,
		db:      db,
		metrics: metrics,
	}
}

// parseParams parses experiment parameters from JSON string
func parseParams(experiment *storage.Experiment) (map[string]string, error) {
	if experiment == nil {
		return nil, fmt.Errorf("experiment cannot be nil")
	}
	
	var params map[string]string
	if experiment.Parameters == "" {
		params = make(map[string]string)
	} else if err := json.Unmarshal([]byte(experiment.Parameters), &params); err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}
	return params, nil
}

// getExperimentParams extracts common parameters for experiments
func getExperimentParams(params map[string]string, paramName string, defaultValue int, unit string) (*ExperimentParams, error) {
	// Check for required parameters
	namespace, ok := params["namespace"]
	if !ok || namespace == "" {
		return nil, fmt.Errorf("missing required parameter: namespace")
	}
	
	// Get selector parameter
	selector, ok := params["selector"]
	if !ok || selector == "" {
		return nil, fmt.Errorf("missing required parameter: selector")
	}
	
	// Get value parameter with default
	valueStr, ok := params[paramName]
	value := defaultValue
	if ok && valueStr != "" {
		parsedValue, err := strconv.Atoi(valueStr)
		if err != nil {
			log.Printf("Invalid %s parameter: %v, using default", paramName, err)
		} else if parsedValue > 0 {
			value = parsedValue
		}
	}
	
	// Log the value being used
	log.Printf("Using %s of %d%s for experiment", paramName, value, unit)
	
	return &ExperimentParams{
		Namespace: namespace,
		Selector:  selector,
		Value:     value,
		ValueUnit: unit,
	}, nil
}

// executeGenericExperiment executes a generic experiment with common behavior
func (e *Executor) executeGenericExperiment(ctx context.Context, experiment *storage.Experiment, experimentType string, params *ExperimentParams) (*experiments.ExperimentResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if experiment == nil {
		return nil, fmt.Errorf("experiment cannot be nil")
	}
	
	// Create a result object
	result := &experiments.ExperimentResult{
		ExperimentType: experimentType,
		StartTime:      time.Now(),
		Success:        true,
	}
	
	// Simulate experiment execution
	select {
	case <-time.After(time.Duration(experiment.Duration) * time.Second):
		// Experiment completed successfully
	case <-ctx.Done():
		// Experiment was cancelled
		result.Error = "Experiment cancelled"
		result.Success = false
		return result, ctx.Err()
	}
	
	result.EndTime = time.Now()
	return result, nil
}

// ExecuteExperiment executes a chaos experiment identified by experimentID and returns the result
func (e *Executor) ExecuteExperiment(experimentID string) (*experiments.ExperimentResult, error) {
	// Validate experiment ID
	if experimentID == "" {
		return nil, fmt.Errorf("experiment ID cannot be empty")
	}

	// Get the experiment from the database
	experiment, err := e.db.GetExperiment(experimentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get experiment: %w", err)
	}

	// Update experiment status to running
	if err := e.db.UpdateExperimentStatus(experimentID, storage.StatusRunning); err != nil {
		return nil, fmt.Errorf("failed to update experiment status: %w", err)
	}

	// Increment active experiments metric
	e.metrics.ActiveExperiments.Inc()
	defer e.metrics.ActiveExperiments.Dec()

	// Validate experiment duration
	if experiment.Duration <= 0 {
		return nil, fmt.Errorf("invalid experiment duration: %d, must be greater than 0", experiment.Duration)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(experiment.Duration)*time.Second)
	defer cancel()

	// Execute the experiment based on its type
	var result *experiments.ExperimentResult
	var execErr error

	// Validate experiment type
	if experiment.Type == "" {
		execErr = fmt.Errorf("experiment type cannot be empty")
	} else {
		// Use a map to simplify experiment type selection
		executors := map[storage.ExperimentType]func(context.Context, *storage.Experiment) (*experiments.ExperimentResult, error){
			storage.PodFailure:   e.executePodFailure,
			storage.NetworkDelay: e.executeNetworkDelay,
			storage.CPUStress:    e.executeCPUStress,
			storage.MemoryStress: e.executeMemoryStress,
		}
		
		executor, exists := executors[experiment.Type]
		if exists {
			result, execErr = executor(ctx, experiment)
		} else {
			execErr = fmt.Errorf("unsupported experiment type: %s", experiment.Type)
		}
	}

	// Update experiment status based on result
	var status storage.ExperimentStatus
	if execErr != nil {
		status = storage.StatusFailed
		e.metrics.ExperimentsFailed.Inc()
	} else {
		status = storage.StatusCompleted
		e.metrics.ExperimentsSucceeded.Inc()
	}

	if err := e.db.UpdateExperimentStatus(experimentID, status); err != nil {
		log.Printf("Failed to update experiment status: %v", err)
	}

	// Save the result
	if result != nil {
		result.ExperimentID = experimentID
		result.ID = uuid.New().String()
		result.CalculateDuration()
		
		// Record experiment duration in metrics
		e.metrics.ExperimentDuration.Observe(result.Duration)
	}

	return result, execErr
}

// executePodFailure executes a pod failure experiment
func (e *Executor) executePodFailure(ctx context.Context, experiment *storage.Experiment) (*experiments.ExperimentResult, error) {
	// Parse parameters
	params, err := parseParams(experiment)
	if err != nil {
		return nil, err
	}
	
	// Get common parameters
	experimentParams, err := getExperimentParams(params, "percentage", 100, "%")
	if err != nil {
		return nil, err
	}

	// Create and run the experiment
	podFailure := experiments.NewPodFailureExperiment(
		e.client.GetClientset(),
		experimentParams.Namespace,
		experimentParams.Selector,
		experiment.Duration,
		experimentParams.Value,
	)

	return podFailure.Run(ctx)
}

// executeNetworkDelay executes a network delay experiment
func (e *Executor) executeNetworkDelay(ctx context.Context, experiment *storage.Experiment) (*experiments.ExperimentResult, error) {
	// Parse parameters
	params, err := parseParams(experiment)
	if err != nil {
		return nil, err
	}
	
	// Get common parameters
	experimentParams, err := getExperimentParams(params, "delay", 100, "ms")
	if err != nil {
		return nil, err
	}
	
	// Use generic experiment execution
	return e.executeGenericExperiment(ctx, experiment, "network-delay", experimentParams)
}

// executeCPUStress executes a CPU stress experiment
func (e *Executor) executeCPUStress(ctx context.Context, experiment *storage.Experiment) (*experiments.ExperimentResult, error) {
	// Parse parameters
	params, err := parseParams(experiment)
	if err != nil {
		return nil, err
	}
	
	// Get common parameters
	experimentParams, err := getExperimentParams(params, "load", 80, "%")
	if err != nil {
		return nil, err
	}
	
	// Use generic experiment execution
	return e.executeGenericExperiment(ctx, experiment, "cpu-stress", experimentParams)
}

// executeMemoryStress executes a memory stress experiment
func (e *Executor) executeMemoryStress(ctx context.Context, experiment *storage.Experiment) (*experiments.ExperimentResult, error) {
	// Parse parameters
	params, err := parseParams(experiment)
	if err != nil {
		return nil, err
	}
	
	// Get common parameters
	experimentParams, err := getExperimentParams(params, "size", 256, "MB")
	if err != nil {
		return nil, err
	}
	
	// Use generic experiment execution
	return e.executeGenericExperiment(ctx, experiment, "memory-stress", experimentParams)
}