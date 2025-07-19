# Architecture Overview

This document provides a detailed overview of the Chaos Engineering as a Platform architecture.

## System Architecture

The platform follows a layered microservices architecture built on Kubernetes:

![Architecture Diagram](images/architecture.png)

### Layers

#### 1. User Interface Layer

- **Web Dashboard**: Bootstrap 5 and vanilla JavaScript web application for managing experiments
- **CLI Tool**: Command-line interface for automation and scripting
- **API Clients**: Client libraries for various programming languages

#### 2. API Gateway Layer

- **Authentication**: JWT-based authentication and authorization
- **Rate Limiting**: Prevents API abuse
- **Request Routing**: Routes requests to appropriate services

#### 3. Core Services Layer

- **API Server**: RESTful API for managing experiments and targets
- **Chaos Operator**: Kubernetes operator for executing chaos experiments
- **Experiment Scheduler**: Schedules and manages experiment execution
- **Safety System**: Enforces safety guardrails and prevents cascading failures

#### 4. Infrastructure Layer

- **Kubernetes Cluster**: Underlying infrastructure for running the platform
- **PostgreSQL Database**: Stores experiment definitions and results
- **Redis Cache**: Caches frequently accessed data
- **Container Registry**: Stores container images for chaos experiments

#### 5. Monitoring Layer

- **Prometheus**: Collects and stores metrics
- **Grafana**: Visualizes metrics and provides dashboards
- **Alert Manager**: Sends alerts based on predefined conditions

#### 6. External Integrations

- **Cloud Providers**: AWS, GCP, Azure integrations
- **Notification Systems**: Slack, Email, PagerDuty integrations
- **CI/CD Systems**: Jenkins, GitHub Actions, GitLab CI integrations

## Component Interactions

### Experiment Execution Flow

1. User creates an experiment through the UI or API
2. API Server stores the experiment definition in the database
3. User triggers experiment execution
4. API Server updates experiment status to "running"
5. Chaos Operator receives the execution request
6. Operator creates appropriate Kubernetes resources for the experiment
7. Safety System monitors the experiment execution
8. Metrics are collected and stored in Prometheus
9. Experiment completes or is terminated
10. Results are stored in the database
11. Notifications are sent (if configured)

## Key Components

### API Server

The API Server provides a RESTful API for managing experiments and targets. It is implemented in Go using the Gin web framework and communicates with the PostgreSQL database for persistence.

Key responsibilities:
- Experiment and target CRUD operations
- Authentication and authorization
- Input validation
- Database interactions

### Chaos Operator

The Chaos Operator is a Kubernetes operator that executes chaos experiments. It watches for experiment execution requests and creates the necessary Kubernetes resources to inject failures.

Key responsibilities:
- Pod failure experiments
- Network delay experiments
- CPU stress experiments
- Memory stress experiments
- External target experiments

### Safety System

The Safety System enforces safety guardrails to prevent experiments from causing real outages. It monitors experiment execution and can automatically terminate experiments if predefined conditions are met.

Key responsibilities:
- Blast radius control
- Automatic rollback
- Protected namespace enforcement
- Resource utilization monitoring

### Monitoring System

The Monitoring System collects and visualizes metrics from experiments. It uses Prometheus for metric collection and Grafana for visualization.

Key responsibilities:
- Experiment metrics collection
- System metrics collection
- Dashboard visualization
- Alerting

## Data Model

### Experiment

```
Experiment {
  id: UUID
  name: String
  description: String
  type: ExperimentType
  status: ExperimentStatus
  target: String
  parameters: JSON
  created_at: Timestamp
  updated_at: Timestamp
  duration: Integer
}
```

### Target

```
Target {
  id: UUID
  name: String
  description: String
  type: TargetType
  namespace: String
  selector: String
  created_at: Timestamp
}
```

### ExperimentResult

```
ExperimentResult {
  id: UUID
  experiment_id: UUID
  start_time: Timestamp
  end_time: Timestamp
  success: Boolean
  error: String
  metrics: JSON
}
```

## Security Considerations

- **Authentication**: JWT-based authentication for API access
- **Authorization**: Role-based access control for different operations
- **RBAC**: Kubernetes RBAC for operator permissions
- **Network Security**: Service-to-service communication secured with TLS
- **Secrets Management**: Kubernetes Secrets for sensitive information

## Scalability

The platform is designed to scale horizontally:

- **API Server**: Multiple replicas behind a load balancer
- **Chaos Operator**: Multiple replicas with leader election
- **Database**: Connection pooling and potential sharding for large deployments
- **Monitoring**: Federated Prometheus for large-scale metric collection

## Deployment Options

### Single-Cluster Deployment

In this deployment model, the platform runs in the same Kubernetes cluster as the target applications. This is simpler to set up but may not be suitable for production environments where isolation is important.

### Multi-Cluster Deployment

In this deployment model, the platform runs in a dedicated Kubernetes cluster and connects to target clusters using Kubernetes federation or direct API access. This provides better isolation but requires more complex setup.

## Future Architecture Enhancements

- **Multi-tenancy**: Support for multiple teams with isolated environments
- **Chaos Experiment Marketplace**: Predefined experiment templates
- **Machine Learning Integration**: Automated experiment selection based on system behavior
- **Distributed Tracing**: Integration with OpenTelemetry for request tracing
- **Event-Driven Architecture**: Kafka-based event bus for component communication