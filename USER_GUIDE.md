# Chaos Engineering as a Platform - User Guide

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Dashboard Overview](#dashboard-overview)
4. [Creating Targets](#creating-targets)
5. [Creating Experiments](#creating-experiments)
6. [Running Experiments](#running-experiments)
7. [Monitoring Results](#monitoring-results)
8. [API Integration](#api-integration)
9. [Advanced Configuration](#advanced-configuration)
10. [Troubleshooting](#troubleshooting)

## Introduction

Chaos Engineering as a Platform is a comprehensive solution for running controlled chaos experiments in Kubernetes clusters. This guide will help you understand how to use the platform to improve your system's resilience through controlled failure injection.

### What is Chaos Engineering?

Chaos Engineering is the discipline of experimenting on a system to build confidence in its capability to withstand turbulent conditions in production. By deliberately injecting failures, you can identify weaknesses before they manifest in system-wide outages.

### Key Features

- **Pod Failure Experiments**: Terminate pods to test service resilience
- **Network Delay Experiments**: Introduce network latency to test timeout handling
- **CPU Stress Experiments**: Consume CPU resources to test throttling mechanisms
- **Memory Stress Experiments**: Consume memory resources to test OOM handling
- **External Target Experiments**: Test external services through their APIs
- **Experiment Scheduling**: Plan and schedule experiments for automated testing
- **Metrics Collection**: Gather and visualize experiment results
- **API Integration**: Integrate with CI/CD pipelines

## Getting Started

### Prerequisites

Before you begin, ensure you have:

1. A running Kubernetes cluster
2. `kubectl` configured to access your cluster
3. Docker and Docker Compose (for local development)
4. Go 1.19 or later (for development)

### Installation

#### Using Docker Compose (Development)

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/chaos-engineering-as-a-platform.git
   cd chaos-engineering-as-a-platform
   ```

2. Start the local development environment:
   ```bash
   docker-compose up -d
   ```

3. Access the web dashboard at http://localhost

#### Using Kubernetes (Production)

1. Apply the Kubernetes manifests:
   ```bash
   kubectl apply -f deployments/kubernetes/secrets.yaml
   kubectl apply -f deployments/kubernetes/postgres.yaml
   kubectl apply -f deployments/kubernetes/api-server.yaml
   kubectl apply -f deployments/kubernetes/chaos-operator.yaml
   kubectl apply -f deployments/kubernetes/monitoring.yaml
   kubectl apply -f deployments/kubernetes/ingress.yaml
   ```

2. Access the platform using your ingress URL

## Dashboard Overview

The dashboard provides a centralized view of your chaos experiments and their results.

### Dashboard Screenshot

[Dashboard Screenshot]

### Key Components

1. **Summary Cards**: Quick overview of experiment statistics
2. **Recent Experiments**: List of recently executed experiments
3. **Navigation Menu**: Access different sections of the platform
4. **Metrics Visualization**: Charts showing experiment results

## Creating Targets

Targets define the systems that will be subjected to chaos experiments.

### Target Types

1. **Kubernetes Pod**: Target specific pods
2. **Kubernetes Deployment**: Target all pods in a deployment
3. **Kubernetes Service**: Target pods backing a service
4. **External Target**: Target external services via API

### Creating a Target

[Target Creation Screenshot]

1. Navigate to the "Targets" section
2. Click "New Target"
3. Fill in the required information:
   - **Name**: A descriptive name for the target
   - **Type**: Select the target type
   - **Namespace**: For Kubernetes targets, specify the namespace
   - **Selector**: Label selector for Kubernetes targets or URL for external targets
4. Click "Create Target"

## Creating Experiments

Experiments define the chaos conditions to be applied to targets.

### Experiment Types

1. **Pod Failure**: Terminate pods to test resilience
2. **Network Delay**: Introduce network latency
3. **CPU Stress**: Consume CPU resources
4. **Memory Stress**: Consume memory resources
5. **External Target**: Send failure signals to external services

### Creating an Experiment

[Experiment Creation Screenshot]

1. Navigate to the "Experiments" section
2. Click "New Experiment"
3. Fill in the required information:
   - **Name**: A descriptive name for the experiment
   - **Description**: Optional details about the experiment
   - **Type**: Select the experiment type
   - **Target**: Select a previously created target
   - **Duration**: How long the experiment should run (in seconds)
   - **Parameters**: Type-specific parameters (e.g., percentage of pods to fail)
4. Click "Create Experiment"

## Running Experiments

Once created, experiments can be executed manually or scheduled.

### Manual Execution

[Experiment Execution Screenshot]

1. Navigate to the "Experiments" section
2. Find the experiment you want to run
3. Click the "Play" button
4. Monitor the experiment status and results

### Scheduling Experiments

1. Navigate to the "Experiments" section
2. Find the experiment you want to schedule
3. Click "Schedule"
4. Set the schedule parameters (frequency, time window)
5. Click "Save Schedule"

## Monitoring Results

The platform provides detailed monitoring of experiment results.

### Results Dashboard

[Results Dashboard Screenshot]

1. Navigate to the "Results" section
2. View experiment results, including:
   - Success/failure status
   - Duration
   - Affected resources
   - System metrics during the experiment

### Grafana Integration

For more detailed metrics:

1. Click the "Grafana" link in the navigation bar
2. Access pre-configured dashboards for:
   - System metrics during experiments
   - Historical experiment data
   - Target health metrics

## API Integration

The platform provides a RESTful API for integration with CI/CD pipelines.

### API Endpoints

- `POST /api/v1/experiments`: Create a new experiment
- `GET /api/v1/experiments`: List all experiments
- `GET /api/v1/experiments/{id}`: Get an experiment by ID
- `POST /api/v1/experiments/{id}/execute`: Execute an experiment
- `DELETE /api/v1/experiments/{id}`: Delete an experiment
- `GET /api/v1/targets`: List all targets
- `POST /api/v1/targets`: Create a new target

### Example: Creating and Running an Experiment via API

```bash
# Create a target
curl -X POST http://localhost:8080/api/v1/targets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Frontend Service",
    "type": "deployment",
    "namespace": "default",
    "selector": "app=frontend"
  }'

# Create an experiment
curl -X POST http://localhost:8080/api/v1/experiments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Frontend Pod Failure",
    "type": "pod-failure",
    "target": "app=frontend",
    "parameters": {
      "namespace": "default",
      "percentage": "50"
    },
    "duration": 60
  }'

# Execute the experiment
curl -X POST http://localhost:8080/api/v1/experiments/{id}/execute
```

## Advanced Configuration

### Safety Guardrails

Configure safety measures to prevent experiments from causing real outages:

1. Navigate to the "Settings" section
2. Under "Safety Guardrails", configure:
   - **Blast Radius**: Maximum percentage of resources affected
   - **Automatic Rollback**: Conditions for automatic experiment termination
   - **Protected Namespaces**: Namespaces excluded from experiments

### Custom Experiment Templates

Create reusable experiment templates:

1. Navigate to the "Templates" section
2. Click "New Template"
3. Define the template parameters
4. Save the template for future use

## Troubleshooting

### Common Issues

1. **Experiment Fails to Start**
   - Check target exists and is accessible
   - Verify RBAC permissions for the chaos operator

2. **External Target Experiments Fail**
   - Verify the external service URL is correct
   - Check authentication credentials

3. **Metrics Not Showing**
   - Ensure Prometheus and Grafana are running
   - Check service discovery configuration

### Logs

Access logs for troubleshooting:

```bash
# API Server logs
kubectl logs -l app=chaos-api-server

# Chaos Operator logs
kubectl logs -l app=chaos-operator

# Database logs
kubectl logs -l app=postgres
```

For more detailed information, refer to the [API Documentation](API.md) and [Architecture Overview](ARCHITECTURE.md).