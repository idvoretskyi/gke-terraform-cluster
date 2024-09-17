output "cluster_name" {
  description = "The name of the created GKE cluster"
  value       = google_container_cluster.primary.name
}

output "cluster_endpoint" {
  description = "The endpoint of the created GKE cluster"
  value       = google_container_cluster.primary.endpoint
}

output "master_version" {
  description = "The version of Kubernetes used by the GKE cluster"
  value       = google_container_cluster.primary.master_version
}