package operator

// ExperimentController defines the interface for experiment controllers.
// All experiment controllers must implement these methods to be compatible
// with the chaos operator.
type ExperimentController interface {
	// Start starts the experiment
	Start() error
	
	// Stop stops the experiment
	Stop() error
	
	// GetStatus gets the status of the experiment
	GetStatus() string
}