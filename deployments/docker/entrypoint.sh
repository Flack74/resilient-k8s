#!/bin/bash
set -e

# Replace environment variables in kubeconfig
if [ -f "$KUBECONFIG" ]; then
  echo "Setting up Kubernetes configuration..."
  sed -i "s|\${KUBE_TOKEN}|${KUBE_TOKEN}|g" "$KUBECONFIG"
fi

# For development mode, create a mock Kubernetes environment
if [ "$ENVIRONMENT" = "development" ]; then
  echo "Running in development mode with mock Kubernetes client"
  export MOCK_KUBERNETES="true"
fi

# Start the chaos operator
exec /app/chaos-operator