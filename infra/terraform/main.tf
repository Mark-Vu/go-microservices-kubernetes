# Artifact Registry Module
module "artifact_registry" {
  source = "./modules/artifact-registry"

  project_id     = var.project_id
  region         = var.region
  repository_id  = var.artifact_registry_repository
  description    = var.artifact_registry_description
  immutable_tags = false  
}

# IAM Module - Service Account for GitHub Actions
module "iam" {
  source = "./modules/iam"

  project_id          = var.project_id
  service_account_id  = "github-actions-deployer"
  display_name        = "GitHub Actions Deployer"
  description         = "Service account for GitHub Actions CI/CD pipeline"
  repository_location = module.artifact_registry.repository_location
  repository_name     = module.artifact_registry.repository_name

  depends_on = [module.artifact_registry]
}
