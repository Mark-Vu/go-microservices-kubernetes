# Remote state file in Google Cloud Storage
# Bucket name is configured via -backend-config or environment variable
terraform {
  backend "gcs" {
    # bucket = "" is configured via CI/CD workflow
    prefix = "artifact-registry"
  }
}
