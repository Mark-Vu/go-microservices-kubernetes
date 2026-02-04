# Remote state file in Google Cloud Storage
terraform {
  backend "gcs" {
    bucket = "${var.project_id}-tf-state"
    prefix = "artifact-registry"
  }
}
