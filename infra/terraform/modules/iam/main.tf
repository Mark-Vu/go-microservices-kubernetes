# Service account for GitHub Actions
resource "google_service_account" "github_actions" {
  account_id   = var.service_account_id
  display_name = var.display_name
  description  = var.description
}

# Grant Artifact Registry Writer role to service account
resource "google_artifact_registry_repository_iam_member" "writer" {
  project    = var.project_id
  location   = var.repository_location
  repository = var.repository_name
  role       = "roles/artifactregistry.writer"
  member     = "serviceAccount:${google_service_account.github_actions.email}"
}

# Grant project viewer role to service account
resource "google_project_iam_member" "viewer" {
  project = var.project_id
  role    = "roles/viewer"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}
