# Terraform - Artifact Registry

Creates Google Cloud Artifact Registry for Docker images.

## Prerequisites

Create service account manually (one-time):

```bash
# Create service account
gcloud iam service-accounts create ci-cd-admin \
  --display-name="CI/CD Admin"

# Grant permissions
PROJECT_ID=$(gcloud config get-value project)
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:ci-cd-admin@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/editor"

# Create key
gcloud iam service-accounts keys create ci-cd-key.json \
  --iam-account=ci-cd-admin@${PROJECT_ID}.iam.gserviceaccount.com

# Add key contents to GitHub Secrets as GCP_SA_KEY
```

## What This Creates

- Artifact Registry repository: `ride-sharing`

Note: Service account `ci-cd-admin` is created manually (see Prerequisites)

## Quick Start

```bash
# 1. Configure
cd infra/terraform
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your project_id

# 2. Authenticate
gcloud auth application-default login

# 3. Run
terraform init
terraform apply
```

## What Happens

```
terraform apply
  ↓
Creates Artifact Registry in us-west1
  ↓
Outputs registry URL
```

## After Apply

### Get Registry URL

```bash
terraform output artifact_registry_url
# Example: us-west1-docker.pkg.dev/your-project/ride-sharing
```

### Add to GitHub Secrets

Repository → Settings → Secrets → Actions:

- `GCP_SA_KEY`: Contents of `ci-cd-key.json` from Prerequisites
- `GCP_PROJECT_ID`: Your GCP project ID

## Module Structure

```
main.tf
└─ modules/artifact-registry    Creates registry
```

## Common Commands

```bash
terraform plan        Preview changes
terraform apply       Apply changes
terraform destroy     Delete everything
terraform output      Show outputs
```
