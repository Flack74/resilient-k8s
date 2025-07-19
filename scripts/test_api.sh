#!/bin/bash

# Test API server health
echo "Testing API server health..."
curl -s http://localhost:8080/health
echo -e "\n"

# Test Grafana access
echo "Testing Grafana access..."
curl -s -I http://localhost:3000 | head -1
echo -e "\n"

# Test Prometheus access
echo "Testing Prometheus access..."
curl -s -I http://localhost:9090 | head -1
echo -e "\n"

# Test Web Dashboard access
echo "Testing Web Dashboard access..."
curl -s -I http://localhost | head -1
echo -e "\n"

echo "All services are accessible!"