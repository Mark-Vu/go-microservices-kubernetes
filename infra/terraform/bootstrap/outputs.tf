output "state_bucket_name" {
  description = "Name of the Terraform state bucket (use this in backend.tf)"
  value       = google_storage_bucket.tf_state.name
}

output "next_steps" {
  description = "What to do next"
  value       = <<-EOT
    State bucket created: ${google_storage_bucket.tf_state.name}
    
    Next steps:
    1. cd ../terraform
    2. Edit backend.tf - uncomment backend block
    3. Update bucket name to: ${google_storage_bucket.tf_state.name}
    4. Run: terraform init
    5. Your main config will now use remote state
  EOT
}
