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

output "github_actions_service_account_email" {
  description = "GitHub Actions service account email"
  value       = module.iam.service_account_email
}

output "create_key_command" {
  description = "Command to create service account key for GitHub Actions"
  value       = "gcloud iam service-accounts keys create key.json --iam-account=${module.iam.service_account_email}"
}

output "next_steps" {
  description = "Next steps after applying Terraform"
  value       = <<-EOT
    
    ✅ Artifact Registry created successfully!
    
    Next steps:
    
    1. Create service account key:
       gcloud iam service-accounts keys create key.json --iam-account=${module.iam.service_account_email}
    
    2. Add to GitHub Secrets:
       - GCP_SA_KEY: (contents of key.json)
       - GCP_PROJECT_ID: ${var.project_id}
    
    3. Use this image URL in your deployments:
       ${module.artifact_registry.repository_url}/IMAGE_NAME:TAG
       
       Example:
       ${module.artifact_registry.repository_url}/api-gateway:latest
       ${module.artifact_registry.repository_url}/trip-service:latest
    
    4. Push your first image:
       docker tag api-gateway ${module.artifact_registry.repository_url}/api-gateway:latest
       docker push ${module.artifact_registry.repository_url}/api-gateway:latest
  EOT
}
