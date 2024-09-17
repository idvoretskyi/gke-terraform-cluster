provider "google" {
  project = var.project_id
  region  = var.region
}

provider "random" {
  # Random provider to generate a suffix for cluster name
}

# Generate a random suffix for the cluster name
resource "random_string" "random_suffix" {
  length  = 6
  upper   = false
  lower   = true
  numeric = true
  special = false
}

resource "google_container_cluster" "primary" {
  name     = "${var.cluster_name}-${random_string.random_suffix.result}"
  location = var.region

  # Use release channel to automatically use the most recent Kubernetes version
  release_channel {
    channel = "RAPID"  # You can switch to "STABLE" if preferred
  }

  remove_default_node_pool = true
  initial_node_count       = 1

  # Enable autoscaling at the cluster level
  cluster_autoscaling {
    enabled = true
  }

  # Define a preemptible node pool with autoscaling
  node_pool {
    name       = "preemptible-pool"
    node_count = var.min_nodes

    autoscaling {
      min_node_count = var.min_nodes
      max_node_count = var.max_nodes
    }

    node_config {
      machine_type = var.machine_type
      preemptible  = true
      disk_size_gb = 30

      metadata = {
        disable-legacy-endpoints = "true"
      }

      oauth_scopes = [
        "https://www.googleapis.com/auth/cloud-platform",
      ]
    }

    management {
      auto_repair  = true
      auto_upgrade = true
    }
  }
}

output "endpoint" {
  value = google_container_cluster.primary.endpoint
}

output "master_version" {
  value = google_container_cluster.primary.master_version
}

output "cluster_name" {
  value = google_container_cluster.primary.name
}