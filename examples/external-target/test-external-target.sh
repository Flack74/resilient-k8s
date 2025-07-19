#!/bin/bash

# This script demonstrates how to use external targets with the Chaos Engineering Platform

# Start the external service in the background
echo "Starting external service..."
go run external-service.go &
EXTERNAL_PID=$!

# Wait for the service to start
sleep 2

# Create an external target
echo "Creating external target..."
TARGET_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/targets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Local External Service",
    "description": "Local external service for testing",
    "type": "external",
    "namespace": "default",
    "selector": "http://localhost:8090/chaos"
  }')

echo "Target created: $TARGET_RESPONSE"
TARGET_ID=$(echo $TARGET_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

# Create an experiment
echo "Creating experiment..."
EXPERIMENT_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/experiments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "External Service Test",
    "description": "Test external service failure injection",
    "type": "external-target",
    "target": "http://localhost:8090/chaos",
    "parameters": {
      "target_type": "external",
      "endpoint": "/inject-failure",
      "type": "service-failure",
      "cleanup_endpoint": "/reset"
    },
    "duration": 30
  }')

echo "Experiment created: $EXPERIMENT_RESPONSE"
EXPERIMENT_ID=$(echo $EXPERIMENT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

# Test the API before chaos
echo "Testing API before chaos..."
curl -s http://localhost:8090/api | jq .

# Execute the experiment
echo "Executing experiment..."
curl -s -X POST http://localhost:8080/api/v1/experiments/$EXPERIMENT_ID/execute

# Test the API during chaos
echo "Testing API during chaos (should fail)..."
sleep 2
curl -s http://localhost:8090/api | jq .

# Wait for experiment to complete
echo "Waiting for experiment to complete..."
sleep 30

# Test the API after chaos
echo "Testing API after chaos (should work again)..."
curl -s http://localhost:8090/api | jq .

# Clean up
echo "Cleaning up..."
kill $EXTERNAL_PID

echo "Test complete!"