variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "cluster_name" {
  description = "Name of the GKE cluster"
  type        = string
  default     = "ride-sharing"
}

variable "zone" {
  description = "GCP zone for the cluster"
  type        = string
  default     = "us-west1-a"
}
