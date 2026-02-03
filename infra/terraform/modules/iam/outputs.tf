output "service_account_email" {
  description = "Email of the service account"
  value       = google_service_account.github_actions.email
}

output "service_account_id" {
  description = "ID of the service account"
  value       = google_service_account.github_actions.account_id
}

output "service_account_name" {
  description = "Fully qualified name of the service account"
  value       = google_service_account.github_actions.name
}
