# Microservices Ride-Sharing Platform

A backend microservices system for a Uber-style ride-sharing app built with Go, Docker, and Kubernetes.

## What This Is

Event-driven microservices for ride-sharing (trips, drivers, payments). Deployed to GKE with automated CI/CD pipeline.

## Clean Architecture

### Folder Structure
```
services/trip-service/
├── cmd/
│   └── main.go              # Entry point: wires everything together
├── internal/
│   ├── domain/              # Business logic layer (entities + interfaces)
│   │   ├── trip.go          # Core models + repository/service interfaces
│   │   └── ride_fare.go
│   ├── service/             # Application layer (use cases)
│   │   └── service.go       # Business logic implementation
│   └── infrastructure/      # External concerns (DB, HTTP, gRPC)
│       ├── http/            # HTTP handlers
│       ├── repository/      # Database implementations
│       └── grpc/            # gRPC implementations
```
### Why Clean Architecture?
- **Testable**: Mock any layer easily
- **Maintainable**: Clear separation of concerns
- **Flexible**: Swap DB/HTTP without touching business logic
- **Scalable**: Add new features without breaking existing code

## Technical Choices

### Backend
- **Go 1.23+**: Native concurrency (goroutines), fast compilation
- **Clean Architecture**: Domain/Service/Infrastructure layers, testable
- **Graceful Shutdown**: Handles SIGTERM, waits for in-flight requests/messages

### Infrastructure
- **Tilt**: Local development with live reload and auto Kubernetes deployment
- **Kubernetes**: Container orchestration for scaling and self-healing
- **Docker**: Consistent containerized environments across all stages

### Communication
- **gRPC**: High-performance RPC between services
- **WebSockets**: Real-time updates for driver locations and trip status
- **RabbitMQ**: Event-driven async messaging with topic exchanges and fair dispatch

### Observability
- **Jaeger**: Distributed tracing across microservices
- **Structured Logging**: Context propagation for debugging
- **Health Checks**: Kubernetes liveness and readiness probes

### Cloud & Deployment
- **Terraform**: Infrastructure as code (GKE + Artifact Registry)
- **GitHub Actions**: Automated CI/CD on every push
- **ArgoCD**: GitOps - cluster syncs with Git automatically
- **GKE Autopilot**: Managed Kubernetes, auto-scales nodes

## Trip Scheduling Flow
[![](https://mermaid.ink/img/pako:eNqNVt9v2jAQ_lcsP21qGvGjZSEPlSpaTX1YxWDVpAmpMvZBIkicOQ6UVf3fd4mdEgcKzQOK47v7vjt_d-aVcimAhnSW5vC3gJTDXcyWiiWzlOCTMaVjHmcs1eQpB3X49Xb88J1p2LLd4d4vFWdTUJuYw-HmnYo3oM5sH34fs10Cqf7Qb6oRFL-bnZLz5c3NnmRIRgrwteJGJmXOuTa2eyP0aFAPyXIyHtWO5YaxN7-PEoNJpNrM1nOSCw3Y_QuPWLq0nBvWl4h30fIos_Bhg5n6vMIVDTgVLyNN5II4LO9L65A8wrbyJimAyAkjwlbSBHBwENisQ2vl80T4pfezapbGRW1RHckkYaloILMNi9dsvgYXtHUQv2E-lXwF-hCccQ5ZE3sNiwZ0A_O2sjSw6oPTbB_nWbT3TB3dGEiykGrLlABBtCQTNp_H-sdP8sUwK3EmkGcS2-lnAQV8rUtwVCchGSvJIc8tJ2KoMGxDjxSZKIVapZZrpovc9_2j4nF7whGPifvM8jxepp8XkW2SzAQmOVKMZXoyF6_Nwq5bunetkLxp2HfIUQR8JQts5Bqz9DJGx3K1ZtZdHAVpCc9mZStkc3s-0WZtTFukOsGnB5QeE7u6PM4gKScQsozktmF_TmuN1qidUHZJa6q9V04m2RowlrV1SnbQc5GUq33YacFL_R2faHtPzxXt0ZN1O-6iXbRW1Zu4HxeiqjTJivk6ziPTcp-U1aXD-NPgZ46a21KL061Qj6kzc38_facarzCyv1tOdGgVk7OUpCipOSxjZEI9moBKWCzwKn8tQ8yojiCBGQ3xVTC1muEV_4Z2rNByuks5DbUqwKNKFsuIhgu2znFlZo79C1Cb4L36R8rmkoav9IWGvW_-1XVn0O_1-kE3GAyHgUd3-Lnb8fu9frc_xKfbvQ6CN4_-qyJ0_KDX7Q86QTDoDAfD66ve23_1IPGQ?type=png)](https://mermaid.live/edit#pako:eNqNVt9v2jAQ_lcsP21qGvGjZSEPlSpaTX1YxWDVpAmpMvZBIkicOQ6UVf3fd4mdEgcKzQOK47v7vjt_d-aVcimAhnSW5vC3gJTDXcyWiiWzlOCTMaVjHmcs1eQpB3X49Xb88J1p2LLd4d4vFWdTUJuYw-HmnYo3oM5sH34fs10Cqf7Qb6oRFL-bnZLz5c3NnmRIRgrwteJGJmXOuTa2eyP0aFAPyXIyHtWO5YaxN7-PEoNJpNrM1nOSCw3Y_QuPWLq0nBvWl4h30fIos_Bhg5n6vMIVDTgVLyNN5II4LO9L65A8wrbyJimAyAkjwlbSBHBwENisQ2vl80T4pfezapbGRW1RHckkYaloILMNi9dsvgYXtHUQv2E-lXwF-hCccQ5ZE3sNiwZ0A_O2sjSw6oPTbB_nWbT3TB3dGEiykGrLlABBtCQTNp_H-sdP8sUwK3EmkGcS2-lnAQV8rUtwVCchGSvJIc8tJ2KoMGxDjxSZKIVapZZrpovc9_2j4nF7whGPifvM8jxepp8XkW2SzAQmOVKMZXoyF6_Nwq5bunetkLxp2HfIUQR8JQts5Bqz9DJGx3K1ZtZdHAVpCc9mZStkc3s-0WZtTFukOsGnB5QeE7u6PM4gKScQsozktmF_TmuN1qidUHZJa6q9V04m2RowlrV1SnbQc5GUq33YacFL_R2faHtPzxXt0ZN1O-6iXbRW1Zu4HxeiqjTJivk6ziPTcp-U1aXD-NPgZ46a21KL061Qj6kzc38_facarzCyv1tOdGgVk7OUpCipOSxjZEI9moBKWCzwKn8tQ8yojiCBGQ3xVTC1muEV_4Z2rNByuks5DbUqwKNKFsuIhgu2znFlZo79C1Cb4L36R8rmkoav9IWGvW_-1XVn0O_1-kE3GAyHgUd3-Lnb8fu9frc_xKfbvQ6CN4_-qyJ0_KDX7Q86QTDoDAfD66ve23_1IPGQ)


## Local Development

**Prerequisites:**
- Docker Desktop (includes Kubernetes)
- Go 1.23+
- Tilt
- kubectl

**Install (macOS):**
```bash
brew install go tilt kubectl
# Install Docker Desktop manually from docker.com
```

**Run:**
```bash
tilt up
```

**Monitor:**
```bash
kubectl get pods
# or
minikube dashboard
```

## Deployment

Fully automated CI/CD pipeline using Terraform, GitHub Actions, and ArgoCD.

### Infrastructure (Terraform)

Provisions GKE cluster + Artifact Registry.

```bash
cd infra/terraform

# Configure
cp terraform.tfvars.example terraform.tfvars
# Edit: set your project_id

# Apply
gcloud auth application-default login
terraform init
terraform apply
```

**What gets created:**
- GKE cluster (autopilot, 3 zones)
- Artifact Registry (`ride-sharing` repo)
- IAM permissions

### CI/CD Pipeline (GitHub Actions)

Automated on every push to `main`:

**1. Build & Push**
- Compiles Go services
- Builds Docker images (tagged with commit SHA)
- Pushes to Artifact Registry

**2. Update Manifests**
- Updates `kustomization.yaml` with new image tag
- Commits back to repo (GitOps)

**Setup GitHub Secrets:**
```
GCP_SA_KEY       → Service account JSON (from Terraform output)
GCP_PROJECT_ID   → Your GCP project ID
GCP_REGION       → us-west1
```

### GitOps (ArgoCD)

Auto-syncs cluster with Git repo.

**Install ArgoCD:**
```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Access UI
kubectl port-forward svc/argocd-server -n argocd 8080:443
# Username: admin
# Password: kubectl -n argocd get secret argocd-initial-admin-password -o jsonpath="{.data.password}" | base64 -d
```

**Deploy app:**
```bash
kubectl apply -f infra/argocd/application.yaml
```

**What happens:**
- ArgoCD watches `infra/production/k8s/overlays/production`
- Detects changes (from GitHub Actions commit)
- Auto-syncs to cluster
- Self-heals if manual changes occur

### Deployment Flow

```
Push to main
  ↓
GitHub Actions builds images
  ↓
Updates kustomization.yaml (new image tag)
  ↓
Commits back to repo
  ↓
ArgoCD detects change
  ↓
Syncs to GKE cluster
  ↓
Pods restart with new image
```

**No manual kubectl commands needed.**

### Check Deployment

```bash
# Get external IP
kubectl get ingress

# Check ArgoCD sync status
kubectl get applications -n argocd

# View pods
kubectl get pods

# Switch back to local cluster
kubectl config use-context docker-desktop  # or minikube
```

### Access Production

ArgoCD creates Ingress with Google-managed SSL.

```bash
kubectl get ingress
# Visit: https://<INGRESS_IP>
```
