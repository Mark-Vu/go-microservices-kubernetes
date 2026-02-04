variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCS bucket location"
  type        = string
  default     = "us-west1"
}
