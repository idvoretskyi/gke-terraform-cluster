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

  # Set deletion protection to false to allow deletion of the cluster
  deletion_protection = false

  release_channel {
    channel = "RAPID"  # You can switch to "STABLE" if preferred
  }

  remove_default_node_pool = true

  # Enable autoscaling at the cluster level
  cluster_autoscaling {
    enabled = true

    # Specify resource limits for CPU and memory (adjust as needed)
    resource_limits {
      resource_type = "cpu"
      minimum = 1
      maximum = 100
    }

    resource_limits {
      resource_type = "memory"
      minimum = 1
      maximum = 200
    }
  }

  # Define a preemptible node pool with autoscaling (no custom name)
  node_pool {
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