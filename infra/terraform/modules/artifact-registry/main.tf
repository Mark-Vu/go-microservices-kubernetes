# Enable Artifact Registry API
resource "google_project_service" "artifact_registry" {
  project = var.project_id
  service = "artifactregistry.googleapis.com"

  disable_on_destroy = false
}

# Create Artifact Registry repository
resource "google_artifact_registry_repository" "main" {
  location      = var.region
  repository_id = var.repository_id
  description   = var.description
  format        = "DOCKER"

  docker_config {
    immutable_tags = var.immutable_tags
  }
  vulnerability_scanning_config {
    enablement_config = "DISABLED" # Possible values: INHERITED, DISABLED
  }

  depends_on = [google_project_service.artifact_registry]
}
