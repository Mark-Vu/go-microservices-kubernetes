variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP region for the repository"
  type        = string
}

variable "repository_id" {
  description = "The ID of the repository"
  type        = string
}

variable "description" {
  description = "Description of the repository"
  type        = string
  default     = "Docker repository for container images"
}

variable "immutable_tags" {
  description = "Whether tags are immutable"
  type        = bool
  default     = false
}
