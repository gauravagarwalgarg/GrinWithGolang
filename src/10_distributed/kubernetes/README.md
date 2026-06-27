# Kubernetes Deployment Manifests

> Deploy the clean backend to Kubernetes with proper resource management, scaling, and configuration.

## Manifests

| File | What It Does |
|------|-------------|
| `namespace.yaml` | Isolated namespace for the app |
| `configmap.yaml` | Non-secret configuration (DB host, timeouts) |
| `secret.yaml` | Sensitive values (tokens, passwords) base64 encoded |
| `deployment.yaml` | App pods: replicas, resource limits, health checks |
| `service.yaml` | ClusterIP service (internal load balancer) |
| `hpa.yaml` | Horizontal Pod Autoscaler (scale based on CPU/memory) |
| `ingress.yaml` | External access with TLS termination |

## Quick Start (Minikube)

```bash
# Start minikube
minikube start

# Apply all manifests
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f hpa.yaml

# Check status
kubectl -n clean-backend get all

# Port-forward to test locally
kubectl -n clean-backend port-forward svc/clean-backend 8080:8080
```

## Architecture on K8s

```
Internet → Ingress → Service → Pod(s) → App Container
                                  ↓
                              MongoDB (StatefulSet or managed)
```
