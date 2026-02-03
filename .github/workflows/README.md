# GitHub Actions CI/CD Setup

## Overview

These workflows automatically build, push, and deploy your services to Google Kubernetes Engine (GKE) using a GitOps approach.

## Workflows

- **`build-and-deploy.yml`** - Reusable workflow for building and deploying services
- **`api-gateway.yml`** - Triggers build and deploy for API Gateway service
- **`trip-service.yml`** - Triggers build and deploy for Trip Service

## How It Works

### Trigger Flow
```
Code Change (push to main)
  ↓
GitHub Actions builds Docker image
  ↓
Pushes to Google Artifact Registry
  ↓
Updates k8s manifest with new image tag (commit SHA)
  ↓
Commits manifest back to repo
  ↓
ArgoCD detects change and deploys to GKE
```

### Architecture

**Reusable Workflow Pattern:**
```
api-gateway.yml (trigger) ──┐
                            ├──> build-and-deploy.yml (reusable logic)
trip-service.yml (trigger) ─┘
```

### Path Filters

Each workflow only triggers when relevant files change:

| Workflow | Triggers On |
|----------|-------------|
| `api-gateway.yml` | `services/api-gateway/**`, `shared/**` |
| `trip-service.yml` | `services/trip-service/**`, `shared/**`, `proto/**` |

### Adding New Services

To add a new service, create a workflow file:

```yaml
name: Build and Deploy New Service

on:
  push:
    branches: [main]
    paths:
      - 'services/new-service/**'

jobs:
  deploy:
    uses: ./.github/workflows/build-and-deploy.yml
    with:
      service_name: new-service
      dockerfile_path: infra/production/docker/new-service.Dockerfile
      manifest_path: infra/production/k8s/new-service-deployment.yaml
      go_version: '1.21'
    secrets:
      gcp_sa_key: ${{ secrets.GCP_SA_KEY }}
      gcp_project_id: ${{ secrets.GCP_PROJECT_ID }}
```

## Prerequisites

### 1. Google Cloud Setup

#### Create Artifact Registry Repository
```bash
# Enable Artifact Registry API
gcloud services enable artifactregistry.googleapis.com

# Create repository
gcloud artifacts repositories create ride-sharing \
  --repository-format=docker \
  --location=us-central1 \
  --description="Ride sharing application images"
```

#### Create Service Account
```bash
# Create service account
gcloud iam service-accounts create github-actions-deployer \
  --description="Service account for GitHub Actions" \
  --display-name="GitHub Actions Deployer"

# Grant permissions
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:github-actions-deployer@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

# Generate and download key
gcloud iam service-accounts keys create key.json \
  --iam-account=github-actions-deployer@YOUR_PROJECT_ID.iam.gserviceaccount.com

# View the key (copy this for GitHub Secrets)
cat key.json
```

### 2. GitHub Secrets Setup

Go to your GitHub repository → Settings → Secrets and variables → Actions → New repository secret

Add these secrets:

| Secret Name | Value | How to Get |
|-------------|-------|------------|
| `GCP_SA_KEY` | Service account JSON key | Output from `cat key.json` above |
| `GCP_PROJECT_ID` | Your GCP project ID | Run `gcloud config get-value project` |

**Example:**
```bash
# Get your project ID
gcloud config get-value project
# Output: my-project-12345

# Copy entire JSON key
cat key.json
# Copy the entire JSON output
```

### 3. Update Image References in Manifests

The workflows expect image names to follow this pattern:
```
us-docker.pkg.dev/YOUR_PROJECT_ID/ride-sharing/IMAGE_NAME:TAG
```

Make sure your k8s deployment manifests use this format:

**Before:**
```yaml
image: ride-sharing/api-gateway
```

**After:**
```yaml
image: us-docker.pkg.dev/YOUR_PROJECT_ID/ride-sharing/api-gateway:latest
```

Run this to update all manifests:
```bash
PROJECT_ID="your-actual-project-id"

sed -i "s|image: ride-sharing/api-gateway|image: us-docker.pkg.dev/$PROJECT_ID/ride-sharing/api-gateway:latest|g" \
  infra/production/k8s/api-gateway-deployment.yaml

sed -i "s|image: ride-sharing/trip-service|image: us-docker.pkg.dev/$PROJECT_ID/ride-sharing/trip-service:latest|g" \
  infra/production/k8s/trip-service-deployment.yaml
```

## Testing the Workflows

### Test Locally (Optional)

You can test Docker builds locally before pushing:

```bash
# Test API Gateway build
docker build -f infra/production/docker/api-gateway.Dockerfile -t test-api-gateway .

# Test Trip Service build
docker build -f infra/production/docker/trip-service.Dockerfile -t test-trip-service .
```

### First Deployment

1. **Push a change** to trigger the workflow:
   ```bash
   # Make a small change
   echo "# CI/CD configured" >> README.md
   git add .
   git commit -m "Setup CI/CD pipeline"
   git push origin main
   ```

2. **Watch the workflow**:
   - Go to GitHub → Actions tab
   - Click on the running workflow
   - Monitor the steps

3. **Verify images**:
   ```bash
   # List images in Artifact Registry
   gcloud artifacts docker images list us-docker.pkg.dev/YOUR_PROJECT_ID/ride-sharing
   ```

## Monitoring

### GitHub Actions
- **Logs**: GitHub → Actions tab → Select workflow run
- **Status**: You'll see ✅ or ❌ on commits

### Artifact Registry
```bash
# List all images
gcloud artifacts docker images list \
  us-docker.pkg.dev/YOUR_PROJECT_ID/ride-sharing

# View image details
gcloud artifacts docker images describe \
  us-docker.pkg.dev/YOUR_PROJECT_ID/ride-sharing/api-gateway:COMMIT_SHA
```

### Check Manifest Updates
After a successful workflow, you should see a new commit from `github-actions[bot]`:
```
Update api-gateway image to abc123def456
```

## Troubleshooting

### Build Fails: "Permission Denied"
**Cause**: Service account doesn't have proper permissions

**Fix**:
```bash
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:github-actions-deployer@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"
```

### Build Fails: "Repository not found"
**Cause**: Artifact Registry repository doesn't exist

**Fix**:
```bash
gcloud artifacts repositories create ride-sharing \
  --repository-format=docker \
  --location=us-central1
```

### Manifest Update Fails: "Nothing to commit"
**Cause**: This is normal! It means the image reference didn't need updating

### Image Tag Pattern Issues
Make sure your sed pattern matches your actual image reference format. Check with:
```bash
grep "image:" infra/production/k8s/*.yaml
```

## Image Tagging Strategy

Each successful build creates TWO tags:

1. **Commit SHA tag** (immutable):
   - Format: `api-gateway:abc123def456`
   - Used in k8s manifests for traceability
   
2. **Latest tag** (moving):
   - Format: `api-gateway:latest`
   - For debugging/testing

## Security Best Practices

✅ **DO:**
- Keep service account keys secure
- Use least-privilege IAM roles
- Rotate service account keys regularly
- Enable vulnerability scanning in Artifact Registry

❌ **DON'T:**
- Commit service account keys to Git
- Use overly permissive IAM roles
- Share GitHub Secrets

## Next Steps

Once workflows are running:
1. Set up ArgoCD to watch your manifests
2. Configure branch protection rules
3. Add PR-based builds (without deploy)
4. Set up Slack/Discord notifications

## Questions?

Check the workflow run logs in GitHub Actions for detailed error messages.
