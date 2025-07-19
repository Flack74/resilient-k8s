# Chaos Engineering as a Platform Helm Chart

This Helm chart deploys the Chaos Engineering as a Platform on a Kubernetes cluster.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- PV provisioner support in the underlying infrastructure (for PostgreSQL persistence)

## Installing the Chart

First, update the dependencies:

```bash
helm dependency update ./deployments/helm/chaos-platform
```

This will download the required charts (PostgreSQL, Prometheus, Grafana) into the `charts` directory.

Then, install the chart with the release name `chaos-platform`:

```bash
helm install chaos-platform ./deployments/helm/chaos-platform
```

## Configuration

The following table lists the configurable parameters of the Chaos Platform chart and their default values.

| Parameter                         | Description                                      | Default                           |
|-----------------------------------|--------------------------------------------------|-----------------------------------|
| `replicaCount`                    | Number of replicas                               | `1`                               |
| `image.repository`                | Image repository                                 | `chaos-engineering-as-a-platform` |
| `image.tag`                       | Image tag                                        | `latest`                          |
| `image.pullPolicy`                | Image pull policy                                | `IfNotPresent`                    |
| `apiServer.port`                  | API server port                                  | `8080`                            |
| `apiServer.resources`             | API server resource requests/limits              | See `values.yaml`                 |
| `chaosOperator.resources`         | Chaos operator resource requests/limits          | See `values.yaml`                 |
| `database.host`                   | PostgreSQL host                                  | `postgres`                        |
| `database.port`                   | PostgreSQL port                                  | `5432`                            |
| `database.name`                   | PostgreSQL database name                         | `chaos_platform`                  |
| `database.user`                   | PostgreSQL user                                  | `postgres`                        |
| `database.existingSecret`         | Existing secret with database credentials        | `chaos-platform-db-credentials`   |
| `monitoring.prometheus.enabled`   | Enable Prometheus integration                    | `true`                            |
| `monitoring.grafana.enabled`      | Enable Grafana integration                       | `true`                            |
| `monitoring.grafana.adminPassword`| Grafana admin password                           | `admin`                           |
| `ingress.enabled`                 | Enable ingress                                   | `false`                           |
| `ingress.className`               | Ingress class name                               | `""`                              |
| `ingress.annotations`             | Ingress annotations                              | `{}`                              |
| `ingress.hosts`                   | Ingress hosts                                    | See `values.yaml`                 |

## Creating the Database Secret

Before installing the chart, create a secret for the database credentials:

```bash
kubectl create secret generic chaos-platform-db-credentials \
  --from-literal=username=postgres \
  --from-literal=password=your-password
```