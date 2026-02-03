variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "service_account_id" {
  description = "The service account ID"
  type        = string
}

variable "display_name" {
  description = "Display name for the service account"
  type        = string
  default     = "Service Account"
}

variable "description" {
  description = "Description of the service account"
  type        = string
  default     = "Service account for CI/CD"
}

variable "repository_location" {
  description = "Artifact Registry repository location"
  type        = string
}

variable "repository_name" {
  description = "Artifact Registry repository name"
  type        = string
}
