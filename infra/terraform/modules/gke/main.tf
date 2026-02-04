# Enable required APIs
resource "google_project_service" "container" {
  project            = var.project_id
  service            = "container.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "compute" {
  project            = var.project_id
  service            = "compute.googleapis.com"
  disable_on_destroy = false
}

  # GKE Autopilot Cluster 
  resource "google_container_cluster" "main" {
    name     = var.cluster_name
    location = var.zone # Zonal cluster
  
    enable_autopilot = true
  
    # Use default VPC
    network    = "default"
    subnetwork = "default"

    release_channel {
      channel = "REGULAR"
    }

  depends_on = [
    google_project_service.container,
    google_project_service.compute,
  ]
}
