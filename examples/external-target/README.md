# External Target Example

This example demonstrates how to use an external target with the Chaos Engineering Platform.

## Creating an External Target

```bash
# Create an external target
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

## Creating an External Experiment

```bash
# Create an experiment targeting the external service
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

## Executing the Experiment

```bash
# Execute the experiment (replace with your experiment ID)
curl -X POST http://localhost:8080/api/v1/experiments/{experiment-id}/execute
```

## How It Works

1. The platform sends a POST request to the external target URL with the specified endpoint
2. The request includes headers with experiment information
3. After the experiment duration, a cleanup request is sent to restore normal operation

## Example External Target Implementation

See the `external-service.go` file in this directory for a simple implementation of an external target service that can receive chaos experiments.