variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP region for resources"
  type        = string
  default     = "us-west1"
}

variable "artifact_registry_repository" {
  description = "Artifact Registry repository name"
  type        = string
  default     = "ride-sharing"
}

variable "artifact_registry_description" {
  description = "Description for the Artifact Registry repository"
  type        = string
  default     = "Docker images for ride sharing microservices"
}
