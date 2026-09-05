# deploy-demo

A demonstration Go service and educational project focused on building a CI/CD pipeline and deploying an application to Kubernetes/k3s.

This repository demonstrates the complete path from source-code checks to Docker image publishing and manual Helm release deployment to a cluster.

## Features

- running a Go application in a container;
- building Docker images with Docker Buildx;
- automated checks with GitHub Actions;
- Helm chart validation with `helm lint` and `helm template`;
- publishing images to Docker Hub;
- deploying the application to Kubernetes/k3s with Helm;
- exposing application metrics and collecting them through a Prometheus ServiceMonitor;
- publishing the service over HTTPS through a Kubernetes Ingress.

## Tech stack

- Go;
- Docker;
- Docker Buildx;
- Kubernetes / k3s;
- Helm;
- GitHub Actions;
- Docker Hub;
- Prometheus;
- Grafana;
- Traefik;
- TLS / cert-manager.

## CI/CD pipeline

The GitHub Actions pipeline is divided into sequential jobs using `needs`:

```text
go test
   ↓
helm lint + helm template
   ↓
Docker build + push
   ↓
manual deployment to k3s
```

### Application checks

The first stage runs the Go application tests:

```bash
go test ./...
```

### Helm chart validation

Before the image is built, the Helm chart is validated for correctness:

```bash
helm lint ./chart
helm template deploy-demo ./chart
```

The rendered manifests are stored as workflow artifacts so they can be reviewed separately from the deployment process.

### Image build and publishing

The Docker image is built with Buildx and published to Docker Hub. The commit SHA is used as the image tag, making it possible to identify the exact source revision running in the environment.

Example image:

```text
<dockerhub-username>/deploy-demo:<commit-sha>
```

### Deployment

Deployment is triggered manually after all checks and image publishing have completed successfully:

```bash
helm install deploy-demo ./helm \
  --namespace deploy-demo \
  --create-namespace \
  --set image.repository=YOUR_DOCKERHUB_USERNAME/deploy-demo
  --set image.tag=<commit-sha>
```

## Project structure

```text
.
├── .github/
│   └── workflows/       # GitHub Actions workflows
├── helm/                # Application Helm chart
│   ├── templates/        # Kubernetes manifests
│   ├── Chart.yaml
│   └── values.yaml
├── Dockerfile            # Application image definition
├── go.mod                # Go dependencies
├── go.sum
└── ...                   # Service source code and configuration
```

> If the directory names in the repository differ, update this section to match the actual project structure.

## Local development

### Run with Go

```bash
go run .
```

### Run tests

```bash
go test ./...
```

### Build the binary

```bash
go build -o deploy-demo .
./deploy-demo
```

### Run with Docker

```bash
docker build -t deploy-demo:local .
docker run --rm -p 8080:8080 deploy-demo:local
```

After startup, the application is available at:

```text
http://localhost:8080
```

Adjust the port if necessary according to the application configuration.

## Deploy to Kubernetes

Install the Helm release in the `deploy-demo` namespace:

```bash
helm upgrade --install deploy-demo ./chart \\
  --namespace deploy-demo \\
  --create-namespace
```

Check the status of the deployed resources:

```bash
kubectl get all -n deploy-demo
kubectl get ingress -n deploy-demo
kubectl get pods -n deploy-demo
```

Render the manifests without installing them:

```bash
helm template deploy-demo ./chart
```

Uninstall the release:

```bash
helm uninstall deploy-demo -n deploy-demo
```

## Configuration

The main application settings are located in `chart/values.yaml`.

Example values that are commonly overridden during installation:

```yaml
image:
  repository: <dockerhub-username>/deploy-demo
  tag: latest

service:
  port: 8080

ingress:
  enabled: true
```

For a specific environment, use a separate values file:

```bash
helm install deploy-demo ./helm \
  --namespace deploy-demo \
  --create-namespace \
  --set image.repository=YOUR_DOCKERHUB_USERNAME/deploy-demo
  --set image.tag=<commit-sha>
```

## Monitoring

The application integrates with Prometheus through a `ServiceMonitor`. This allows the metrics endpoint to be discovered automatically after installing the appropriate Prometheus Operator or kube-prometheus-stack.

Check the monitoring resources:

```bash
kubectl get servicemonitor -n deploy-demo
kubectl get service -n deploy-demo
```

Grafana can use the collected metrics to display the application status, request count, errors, and latency.

## Demo

A running example of the application is available at:

[demo.devgit.net](https://demo.devgit.net/)

## Project goals

This project was created as a practical example of an infrastructure workflow:

- automating application checks;
- packaging a Go service into a Docker image;
- managing Kubernetes resources with Helm;
- delivering artifacts through GitHub Actions;
- performing reproducible deployments by commit SHA;
- integrating monitoring with Prometheus.

