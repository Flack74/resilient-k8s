package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Simple external service that can receive chaos experiments
func main() {
	// Initialize state
	state := &ServiceState{
		FailureActive: false,
		LatencyMs:     0,
		ErrorRate:     0,
	}

	// Set up HTTP handlers
	http.HandleFunc("/chaos/inject-failure", func(w http.ResponseWriter, r *http.Request) {
		handleInjectFailure(w, r, state)
	})
	http.HandleFunc("/chaos/reset", func(w http.ResponseWriter, r *http.Request) {
		handleReset(w, r, state)
	})
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		handleAPI(w, r, state)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handleHealth(w, r, state)
	})

	// Start the server
	log.Println("Starting external service on :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

// ServiceState represents the current state of the service
type ServiceState struct {
	FailureActive bool
	LatencyMs     int
	ErrorRate     int
}

// handleInjectFailure injects a failure into the service
func handleInjectFailure(w http.ResponseWriter, r *http.Request, state *ServiceState) {
	// Log experiment details from headers
	experimentID := r.Header.Get("X-Chaos-Experiment-ID")
	experimentType := r.Header.Get("X-Chaos-Experiment-Type")
	duration := r.Header.Get("X-Chaos-Experiment-Duration")

	log.Printf("Received chaos experiment: ID=%s, Type=%s, Duration=%s", 
		experimentID, experimentType, duration)

	// Parse request body for failure parameters
	var params struct {
		FailureType string `json:"failure_type"`
		LatencyMs   int    `json:"latency_ms"`
		ErrorRate   int    `json:"error_rate"`
	}

	// Default to service failure if no body provided
	params.FailureType = "service-failure"

	// Try to decode the request body
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			log.Printf("Error decoding request body: %v", err)
		}
	}

	// Apply the failure based on type
	switch params.FailureType {
	case "service-failure":
		state.FailureActive = true
		log.Println("Injected service failure")
	case "latency":
		state.LatencyMs = params.LatencyMs
		log.Printf("Injected latency: %dms", params.LatencyMs)
	case "error-rate":
		state.ErrorRate = params.ErrorRate
		log.Printf("Injected error rate: %d%%", params.ErrorRate)
	default:
		log.Printf("Unknown failure type: %s", params.FailureType)
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "failure_injected",
		"type":    params.FailureType,
		"id":      experimentID,
		"message": "Failure has been successfully injected",
	})
}

// handleReset resets the service to normal operation
func handleReset(w http.ResponseWriter, r *http.Request, state *ServiceState) {
	// Get experiment ID from header
	experimentID := r.Header.Get("X-Chaos-Experiment-ID")
	log.Printf("Resetting experiment: %s", experimentID)

	// Reset the service state
	state.FailureActive = false
	state.LatencyMs = 0
	state.ErrorRate = 0

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "reset_complete",
		"id":      experimentID,
		"message": "Service has been reset to normal operation",
	})
}

// handleAPI simulates a normal API endpoint that can be affected by chaos
func handleAPI(w http.ResponseWriter, r *http.Request, state *ServiceState) {
	// Check if service failure is active
	if state.FailureActive {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Service is currently unavailable",
		})
		return
	}

	// Apply latency if configured
	if state.LatencyMs > 0 {
		time.Sleep(time.Duration(state.LatencyMs) * time.Millisecond)
	}

	// Apply error rate if configured
	if state.ErrorRate > 0 {
		// Simple random error based on error rate percentage
		if time.Now().UnixNano()%100 < int64(state.ErrorRate) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Internal server error",
			})
			return
		}
	}

	// Normal response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "API is functioning normally",
		"data": map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
		},
	})
}

// handleHealth provides a health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request, state *ServiceState) {
	// Health check should report actual status
	status := "healthy"
	code := http.StatusOK

	if state.FailureActive {
		status = "unhealthy"
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
		"chaos": map[string]interface{}{
			"failure_active": state.FailureActive,
			"latency_ms":     state.LatencyMs,
			"error_rate":     state.ErrorRate,
		},
	})
}