# API Reference

This document provides detailed information about the Chaos Engineering as a Platform API endpoints.

## Base URL

All API endpoints are relative to the base URL:

```
http://localhost:8080/api/v1
```

For production deployments, replace `localhost:8080` with your actual domain.

## Authentication

API requests require authentication using an API key. Include the API key in the `Authorization` header:

```
Authorization: Bearer YOUR_API_KEY
```

## Experiments

### List Experiments

Retrieves a list of all experiments.

**Request**

```
GET /experiments
```

**Response**

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Frontend Pod Failure",
    "description": "Test frontend resilience",
    "type": "pod-failure",
    "status": "completed",
    "target": "app=frontend",
    "parameters": "{\"namespace\":\"default\",\"percentage\":\"50\"}",
    "created_at": "2023-07-19T10:23:54Z",
    "updated_at": "2023-07-19T10:25:54Z",
    "duration": 60
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "API Network Delay",
    "description": "Test timeout handling",
    "type": "network-delay",
    "status": "pending",
    "target": "app=api",
    "parameters": "{\"namespace\":\"default\",\"delay\":\"200\"}",
    "created_at": "2023-07-19T11:30:00Z",
    "updated_at": "2023-07-19T11:30:00Z",
    "duration": 120
  }
]
```

### Get Experiment

Retrieves a specific experiment by ID.

**Request**

```
GET /experiments/{id}
```

**Response**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Frontend Pod Failure",
  "description": "Test frontend resilience",
  "type": "pod-failure",
  "status": "completed",
  "target": "app=frontend",
  "parameters": "{\"namespace\":\"default\",\"percentage\":\"50\"}",
  "created_at": "2023-07-19T10:23:54Z",
  "updated_at": "2023-07-19T10:25:54Z",
  "duration": 60
}
```

### Create Experiment

Creates a new experiment.

**Request**

```
POST /experiments
```

**Request Body**

```json
{
  "name": "Database Memory Stress",
  "description": "Test database OOM handling",
  "type": "memory-stress",
  "target": "app=database",
  "parameters": {
    "namespace": "default",
    "size": "512"
  },
  "duration": 180
}
```

**Response**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "name": "Database Memory Stress",
  "description": "Test database OOM handling",
  "type": "memory-stress",
  "status": "pending",
  "target": "app=database",
  "parameters": "{\"namespace\":\"default\",\"size\":\"512\"}",
  "created_at": "2023-07-19T14:15:00Z",
  "updated_at": "2023-07-19T14:15:00Z",
  "duration": 180
}
```

### Execute Experiment

Executes an existing experiment.

**Request**

```
POST /experiments/{id}/execute
```

**Response**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "status": "running"
}
```

### Delete Experiment

Deletes an experiment.

**Request**

```
DELETE /experiments/{id}
```

**Response**

```json
{
  "message": "experiment deleted"
}
```

## Targets

### List Targets

Retrieves a list of all targets.

**Request**

```
GET /targets
```

**Response**

```json
[
  {
    "id": "650e8400-e29b-41d4-a716-446655440000",
    "name": "Frontend Service",
    "description": "Frontend web service",
    "type": "deployment",
    "namespace": "default",
    "selector": "app=frontend",
    "created_at": "2023-07-18T09:00:00Z"
  },
  {
    "id": "650e8400-e29b-41d4-a716-446655440001",
    "name": "API Service",
    "description": "Backend API service",
    "type": "service",
    "namespace": "default",
    "selector": "app=api",
    "created_at": "2023-07-18T09:05:00Z"
  }
]
```

### Get Target

Retrieves a specific target by ID.

**Request**

```
GET /targets/{id}
```

**Response**

```json
{
  "id": "650e8400-e29b-41d4-a716-446655440000",
  "name": "Frontend Service",
  "description": "Frontend web service",
  "type": "deployment",
  "namespace": "default",
  "selector": "app=frontend",
  "created_at": "2023-07-18T09:00:00Z"
}
```

### Create Target

Creates a new target.

**Request**

```
POST /targets
```

**Request Body**

```json
{
  "name": "Database Service",
  "description": "PostgreSQL database",
  "type": "statefulset",
  "namespace": "default",
  "selector": "app=database"
}
```

**Response**

```json
{
  "id": "650e8400-e29b-41d4-a716-446655440002",
  "name": "Database Service",
  "description": "PostgreSQL database",
  "type": "statefulset",
  "namespace": "default",
  "selector": "app=database",
  "created_at": "2023-07-19T14:20:00Z"
}
```

### Delete Target

Deletes a target.

**Request**

```
DELETE /targets/{id}
```

**Response**

```json
{
  "message": "target deleted"
}
```

## Error Responses

All API endpoints return appropriate HTTP status codes:

- `200 OK`: Request succeeded
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Missing or invalid API key
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

Error responses include a JSON object with an error message:

```json
{
  "error": "experiment not found"
}
```

## Rate Limiting

API requests are rate-limited to 100 requests per minute per API key. If you exceed this limit, you'll receive a `429 Too Many Requests` response.