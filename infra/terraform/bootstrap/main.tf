# ============================================================================
# Terraform Bootstrap - Run Once
# ============================================================================
#
# Purpose: Creates GCS bucket for Terraform state storage
#
# This is a one-time setup that runs with LOCAL state.
# After creating the bucket, the main terraform config will use REMOTE state.

# Run terraform destroy to delete the state bucket
# ============================================================================

terraform {
  required_version = ">= 1.14.4"
  
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.17.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# State storage bucket
# This bucket will store terraform state for the main infrastructure
resource "google_storage_bucket" "tf_state" {
  name     = "${var.project_id}-tf-state"
  location = var.region

  # Prevent accidental public access
  uniform_bucket_level_access = true

  # Keep history of state changes
  versioning {
    enabled = true
  }
}
