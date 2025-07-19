package experiments

import (
	"context"
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodFailureExperiment represents a pod failure chaos experiment
type PodFailureExperiment struct {
	clientset  kubernetes.Interface
	namespace  string
	selector   string
	duration   int
	percentage int
}

// NewPodFailureExperiment creates a new pod failure experiment
func NewPodFailureExperiment(clientset kubernetes.Interface, namespace, selector string, duration, percentage int) *PodFailureExperiment {
	return &PodFailureExperiment{
		clientset:  clientset,
		namespace:  namespace,
		selector:   selector,
		duration:   duration,
		percentage: percentage,
	}
}

// Run executes the pod failure experiment
func (e *PodFailureExperiment) Run(ctx context.Context) (*ExperimentResult, error) {
	log.Printf("Starting pod failure experiment in namespace %s with selector %s", e.namespace, e.selector)
	
	// Create a result object
	result := &ExperimentResult{
		ExperimentType: "pod-failure",
		StartTime:      time.Now(),
		Success:        false,
	}

	// Get pods matching the selector
	pods, err := e.clientset.CoreV1().Pods(e.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: e.selector,
	})
	if err != nil {
		result.Error = fmt.Sprintf("Failed to list pods: %v", err)
		return result, err
	}

	if len(pods.Items) == 0 {
		result.Error = "No pods found matching the selector"
		return result, fmt.Errorf(result.Error)
	}

	// Calculate how many pods to delete
	count := 1
	if e.percentage > 0 {
		count = (len(pods.Items) * e.percentage) / 100
		if count < 1 {
			count = 1
		}
		if count > len(pods.Items) {
			count = len(pods.Items)
		}
	}

	// Delete the pods
	deletedPods := []string{}
	for i := 0; i < count; i++ {
		podName := pods.Items[i].Name
		log.Printf("Deleting pod %s", podName)
		
		err := e.clientset.CoreV1().Pods(e.namespace).Delete(ctx, podName, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Failed to delete pod %s: %v", podName, err)
			continue
		}
		
		deletedPods = append(deletedPods, podName)
	}

	// Record the affected pods
	result.AffectedResources = deletedPods

	// Wait for the specified duration
	select {
	case <-time.After(time.Duration(e.duration) * time.Second):
		// Experiment completed successfully
	case <-ctx.Done():
		// Experiment was cancelled
		result.Error = "Experiment cancelled"
		return result, ctx.Err()
	}

	// Set end time
	result.EndTime = time.Now()
	
	// Set success if we deleted at least one pod
	if len(deletedPods) > 0 {
		result.Success = true
	}

	return result, nil
}