output "project_id" {
  description = "GCP Project ID"
  value       = var.project_id
}

output "region" {
  description = "GCP Region"
  value       = var.region
}

output "artifact_registry_repository_id" {
  description = "Artifact Registry repository ID"
  value       = module.artifact_registry.repository_id
}

output "artifact_registry_location" {
  description = "Artifact Registry location"
  value       = module.artifact_registry.repository_location
}

output "artifact_registry_url" {
  description = "Full URL to push/pull images"
  value       = module.artifact_registry.repository_url
}

output "service_account_note" {
  description = "Service account information"
  value       = "Service account ci-cd-admin@${var.project_id}.iam.gserviceaccount.com is created manually and used for all CI/CD operations"
}
