# Using External Targets

The Chaos Engineering Platform supports running experiments against external targets outside of your Kubernetes cluster. This guide explains how to configure and use external targets.

## What are External Targets?

External targets are systems or services that are not running in your Kubernetes cluster but can be targeted for chaos experiments. These could be:

- External APIs or web services
- On-premises infrastructure
- Cloud services in different environments
- Third-party services that support chaos testing

## Creating an External Target

To create an external target, use the API to create a target with type `external`:

```bash
curl -X POST http://localhost:8080/api/v1/targets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "External API Service",
    "description": "Third-party API service for testing",
    "type": "external",
    "namespace": "default",
    "selector": "https://api.example.com/chaos"
  }'
```

The `selector` field should contain the base URL of the external target.

## Creating an External Experiment

To create an experiment that targets an external system:

```bash
curl -X POST http://localhost:8080/api/v1/experiments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "External API Failure Test",
    "description": "Test resilience of our system when external API fails",
    "type": "external-target",
    "target": "https://api.example.com/chaos",
    "parameters": {
      "target_type": "external",
      "endpoint": "/inject-failure",
      "auth_token": "your-auth-token",
      "type": "service-failure",
      "cleanup_endpoint": "/reset"
    },
    "duration": 60
  }'
```

## Required Parameters for External Targets

When creating an experiment for an external target, you need to include these parameters:

| Parameter | Description |
|-----------|-------------|
| `target_type` | Must be set to `external` |
| `endpoint` | The endpoint to call on the target URL (e.g., `/inject-failure`) |
| `type` | The type of chaos to inject (e.g., `service-failure`, `latency`) |
| `auth_token` | (Optional) Authentication token for the external service |
| `cleanup_endpoint` | (Optional) Endpoint to call for cleanup after the experiment (defaults to `/cleanup`) |

## How It Works

1. When you execute an experiment against an external target, the platform sends a POST request to the target URL + endpoint.
2. The request includes headers with experiment information:
   - `X-Chaos-Experiment-ID`: The ID of the experiment
   - `X-Chaos-Experiment-Type`: The type of chaos to inject
   - `X-Chaos-Experiment-Duration`: The duration in seconds
   - `Authorization`: Bearer token (if provided)
3. After the experiment duration, or if the experiment is stopped, a cleanup request is sent to the cleanup endpoint.

## Implementing an External Target Service

Your external service needs to implement two endpoints:

1. **Chaos Injection Endpoint**: Receives the chaos request and injects the specified failure
2. **Cleanup Endpoint**: Restores normal operation after the experiment

Example implementation in Go:

```go
package main

import (
	"log"
	"net/http"
	"time"
)

var failureActive bool

func main() {
	http.HandleFunc("/chaos/inject-failure", injectFailure)
	http.HandleFunc("/chaos/cleanup", cleanup)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func injectFailure(w http.ResponseWriter, r *http.Request) {
	// Get experiment details from headers
	experimentID := r.Header.Get("X-Chaos-Experiment-ID")
	experimentType := r.Header.Get("X-Chaos-Experiment-Type")
	
	log.Printf("Received chaos experiment request: %s, type: %s", experimentID, experimentType)
	
	// Activate the failure
	failureActive = true
	
	// Return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "failure_injected"}`))
}

func cleanup(w http.ResponseWriter, r *http.Request) {
	// Get experiment ID from header
	experimentID := r.Header.Get("X-Chaos-Experiment-ID")
	
	log.Printf("Cleaning up experiment: %s", experimentID)
	
	// Deactivate the failure
	failureActive = false
	
	// Return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "cleaned_up"}`))
}
```

## Security Considerations

- Use HTTPS for all external target communications
- Implement proper authentication for your external targets
- Consider network security implications when targeting external systems
- Ensure you have permission to run chaos experiments against the target system