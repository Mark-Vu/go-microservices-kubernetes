# Artifact Registry Module
module "artifact_registry" {
  source = "./modules/artifact-registry"

  project_id     = var.project_id
  region         = var.region
  repository_id  = var.artifact_registry_repository
  description    = var.artifact_registry_description
  immutable_tags = false  
}

# Note: Service account (ci-cd-admin) is created manually
# It has Editor role which includes:
# - Artifact Registry writer (push images)
# - IAM admin (run Terraform)
# - Storage admin (access state bucket)
