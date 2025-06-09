output "cluster_name" {
  description = "The name of the GKE cluster"
  value       = google_container_cluster.primary.name
}

output "cluster_endpoint" {
  description = "The endpoint of the GKE cluster"
  value       = google_container_cluster.primary.endpoint
  sensitive   = true
}

output "cluster_ca_certificate" {
  description = "The cluster CA certificate (base64 encoded)"
  value       = google_container_cluster.primary.master_auth[0].cluster_ca_certificate
  sensitive   = true
}

output "cluster_location" {
  description = "The location (region or zone) of the GKE cluster"
  value       = google_container_cluster.primary.location
}

output "cluster_id" {
  description = "The ID of the GKE cluster"
  value       = google_container_cluster.primary.id
}

output "master_version" {
  description = "The version of Kubernetes used by the GKE cluster"
  value       = google_container_cluster.primary.master_version
}

output "node_pool_instance_group_urls" {
  description = "List of instance group URLs for the node pools"
  value       = google_container_node_pool.primary_nodes.instance_group_urls
}

output "kubeconfig_command" {
  description = "Command to configure kubectl for this cluster"
  value       = "gcloud container clusters get-credentials ${google_container_cluster.primary.name} --region ${google_container_cluster.primary.location} --project ${local.project_id}"
}

output "project_id" {
  description = "The GCP project ID being used"
  value       = local.project_id
}