# Sample NGINX Workload

This directory contains the sample NGINX application deployed to validate the GKE cluster functionality.

## Deployment Details

- **Application**: NGINX web server
- **Chart**: bitnami/nginx (version 20.0.8)
- **App Version**: 1.28.0
- **Service Type**: LoadBalancer
- **Replicas**: 2
- **External IP**: http://104.196.125.40

## Deployment Commands

```bash
# Add Bitnami repository
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Deploy NGINX
helm install sample-nginx bitnami/nginx --set service.type=LoadBalancer --set replicaCount=2
```

## Validation

```bash
# Check pods
kubectl get pods -l app.kubernetes.io/name=nginx

# Check service
kubectl get svc sample-nginx

# Test application
curl -I http://104.196.125.40
```

## Status

✅ **Deployment Successful** - The cluster is working correctly and NGINX is responding to HTTP requests.

## Cleanup

To remove the sample application:

```bash
helm uninstall sample-nginx
```