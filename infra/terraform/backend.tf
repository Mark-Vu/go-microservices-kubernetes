# Remote state storage in Google Cloud Storage
#
# Setup (one-time, run manually):
# 
# 1. Create bucket:
#    gcloud storage buckets create gs://PROJECT-ID-tf-state \
#      --location=us-west1 \
#      --uniform-bucket-level-access
#
# 2. Enable versioning:
#    gcloud storage buckets update gs://PROJECT-ID-tf-state --versioning
#
# 3. Uncomment the backend block below and replace PROJECT-ID
#
# 4. Run: terraform init -migrate-state

# terraform {
#   backend "gcs" {
#     bucket = "PROJECT-ID-tf-state"
#     prefix = "artifact-registry"

#   }
# }
