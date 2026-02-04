# Terraform - Artifact Registry

Creates Google Cloud Artifact Registry for Docker images and service account for GitHub Actions.

## What This Creates

- Artifact Registry repository: `ride-sharing`
- Service account: `github-actions-deployer@PROJECT.iam.gserviceaccount.com`
- IAM permissions: Artifact Registry writer + project viewer

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
# Type 'yes' when prompted
```

## What Happens

```
terraform apply
  ↓
Creates Artifact Registry in us-west1
  ↓
Creates service account for GitHub Actions
  ↓
Grants permissions to push Docker images
  ↓
Outputs registry URL
```

## After Apply

### Get Registry URL

```bash
terraform output artifact_registry_url
# Example: us-west1-docker.pkg.dev/your-project/ride-sharing
```

### Create Service Account Key

```bash
SA_EMAIL=$(terraform output -raw github_actions_service_account_email)
gcloud iam service-accounts keys create key.json --iam-account=$SA_EMAIL
```

### Add to GitHub Secrets

Repository → Settings → Secrets → Actions:

- `GCP_SA_KEY`: Contents of `key.json`
- `GCP_PROJECT_ID`: Your GCP project ID

## Module Structure

```
main.tf                      Calls modules
├─ modules/artifact-registry Creates registry
└─ modules/iam              Creates service account
```

## Common Commands

```bash
terraform plan        Preview changes
terraform apply       Apply changes
terraform destroy     Delete everything
terraform output      Show outputs
```