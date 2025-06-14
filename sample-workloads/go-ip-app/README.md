# Go IP Information App

This Go-based web application displays visitor IP information and is designed to be deployed on GKE.

## Features

- Shows visitor's IP address (handles X-Forwarded-For headers)
- Displays server information (hostname, server IP)
- Shows request headers
- Provides both HTML and JSON API endpoints
- Health check endpoint for Kubernetes probes

## Endpoints

- `/` - HTML interface showing IP information
- `/api` - JSON API endpoint
- `/health` - Health check endpoint

## Deployment

### Prerequisites

- GKE cluster running
- Docker configured for GCR
- kubectl configured for your cluster

### Build and Deploy

```bash
# Build and push Docker image
docker build -t gcr.io/PROJECT_ID/go-ip-app:latest .
docker push gcr.io/PROJECT_ID/go-ip-app:latest

# Deploy to Kubernetes
kubectl apply -f deployment.yaml

# Check deployment status
kubectl get pods -l app=go-ip-app
kubectl get svc go-ip-app-service

# Get external IP (may take a few minutes)
kubectl get svc go-ip-app-service -w
```

### Accessing the Application

Once deployed, the LoadBalancer service will provide an external IP. You can access:

- `http://EXTERNAL_IP/` - Web interface
- `http://EXTERNAL_IP/api` - JSON API
- `http://EXTERNAL_IP/health` - Health check

## Configuration

The application uses the following environment variables:

- `PORT` - Server port (default: 8080)

## Docker Image

The application uses a multi-stage Docker build:
1. Build stage: Uses golang:1.21-alpine to compile the application
2. Runtime stage: Uses alpine:latest for a minimal runtime environment

## Kubernetes Resources

- **Deployment**: 2 replicas with resource limits
- **Service**: LoadBalancer type exposing port 80
- **Health Checks**: Liveness and readiness probes on `/health`