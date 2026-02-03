# Remote state storage in Google Cloud Storage
# 
# Before using this backend, create the GCS bucket:
# gcloud storage buckets create gs://YOUR-PROJECT-ID-terraform-state \
#   --location=us-central1 \
#   --uniform-bucket-level-access
#
# Then uncomment the backend block below and update the bucket name

# terraform {
#   backend "gcs" {
#     bucket = "YOUR-PROJECT-ID-terraform-state"
#     prefix = "artifact-registry"
#   }
# }
