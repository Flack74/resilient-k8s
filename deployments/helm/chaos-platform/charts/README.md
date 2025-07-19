# Dependency Charts Directory

This directory is used by Helm to store dependency charts when you run:

```bash
helm dependency update
```

The dependencies are defined in the parent Chart.yaml file and include:

1. PostgreSQL - For database storage
2. Prometheus - For metrics collection
3. Grafana - For metrics visualization

Do not manually add files to this directory as they will be overwritten when updating dependencies.