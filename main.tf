locals {
  project_id   = var.project_id != null && var.project_id != "" ? var.project_id : (
    data.external.gcloud_project.result.project_id != "" ? data.external.gcloud_project.result.project_id : null
  )
  region       = var.region != null && var.region != "" ? var.region : data.external.gcloud_region.result.region
  cluster_name = "${data.external.username.result.username}-${var.cluster_name_suffix}"
}

provider "google" {
  project = local.project_id
  region  = local.region
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
  name     = "${local.cluster_name}-${random_string.random_suffix.result}"
  location = local.region

  # Set deletion protection to false to allow deletion of the cluster
  deletion_protection = false

  release_channel {
    channel = "RAPID"  # Use RAPID for most recent Kubernetes version
  }

  remove_default_node_pool = true

  # Enable cluster autoscaling
  cluster_autoscaling {
    enabled = true
    resource_limits {
      resource_type = "cpu"
      minimum       = 1
      maximum       = 100
    }
    resource_limits {
      resource_type = "memory"
      minimum       = 1
      maximum       = 1000
    }
  }

  # Initial node pool will be removed
  initial_node_count = 1
}

# Create separate cost-optimized node pool
resource "google_container_node_pool" "primary_nodes" {
  name       = "${local.cluster_name}-nodes"
  cluster    = google_container_cluster.primary.id
  location   = local.region
  node_count = var.min_nodes

  autoscaling {
    min_node_count = var.min_nodes
    max_node_count = var.max_nodes
  }

  node_config {
    machine_type = var.machine_type
    spot         = true  # Use spot instances for maximum cost savings
    disk_size_gb = 20    # Reduce disk size for cost savings
    disk_type    = "pd-standard"  # Use standard disks instead of SSD

    metadata = {
      disable-legacy-endpoints = "true"
    }

    # Minimal OAuth scopes for cost optimization
    oauth_scopes = [
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    # Resource constraints for cost control
    resource_labels = {
      environment = "development"
      cost-center = "dev"
    }
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }
}