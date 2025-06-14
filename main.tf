locals {
  project_id = var.project_id != null && var.project_id != "" ? var.project_id : (
    try(data.external.gcloud_project.result.project_id, "") != "" ?
    try(data.external.gcloud_project.result.project_id, "") :
    "your-project-id-here"
  )
  region = var.region != null && var.region != "" ? var.region : (
    try(data.external.gcloud_region.result.region, "") != "" ?
    try(data.external.gcloud_region.result.region, "") :
    "us-central1"
  )
  zone         = "us-central1-a" # Use single zone for cost optimization
  cluster_name = "${try(data.external.username.result.username, "default-user")}-${var.cluster_name_suffix}"
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

# Create VPC network with flow logs enabled
resource "google_compute_network" "vpc" {
  name                    = "${local.cluster_name}-vpc"
  auto_create_subnetworks = false
}

# Create subnet with VPC flow logs and private Google access
resource "google_compute_subnetwork" "subnet" {
  name          = "${local.cluster_name}-subnet"
  ip_cidr_range = "10.10.0.0/24"
  region        = local.region
  network       = google_compute_network.vpc.name

  private_ip_google_access = true

  # Enable VPC Flow Logs and intranode visibility
  log_config {
    aggregation_interval = "INTERVAL_10_MIN"
    flow_sampling        = 0.5
    metadata             = "INCLUDE_ALL_METADATA"
  }

  secondary_ip_range {
    range_name    = "services-range"
    ip_cidr_range = "10.40.0.0/20"
  }

  secondary_ip_range {
    range_name    = "pod-ranges"
    ip_cidr_range = "10.36.0.0/14"
  }
}

# Create firewall rule for the VPI
resource "google_compute_firewall" "allow_internal" {
  name    = "${local.cluster_name}-allow-internal"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "icmp"
  }

  source_ranges = ["10.10.0.0/24", "10.36.0.0/14", "10.40.0.0/20"]
  direction     = "INGRESS"
}

resource "google_container_cluster" "primary" {
  name     = "${local.cluster_name}-${random_string.random_suffix.result}"
  location = local.zone # Use zone instead of region for cost optimization

  # Set deletion protection to false to allow deletion of the cluster
  deletion_protection = false

  # Configure network and subnetwork
  network    = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.subnet.name

  release_channel {
    channel = "RAPID" # Use RAPID for most recent Kubernetes version
  }

  remove_default_node_pool = true

  # Add resource labels for better resource management
  resource_labels = {
    environment = "development"
    cost-center = "dev"
    managed-by  = "terraform"
  }

  # Configure master authorized networks for security
  master_authorized_networks_config {
    cidr_blocks {
      cidr_block   = "0.0.0.0/0"
      display_name = "All networks - development only"
    }
  }

  # Enable private cluster configuration
  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = "172.16.0.32/28"
  }

  # Configure IP allocation policy with alias IP ranges
  ip_allocation_policy {
    cluster_secondary_range_name  = "pod-ranges"
    services_secondary_range_name = "services-range"
  }

  # Enable network policy for enhanced security
  network_policy {
    enabled = true
  }

  # Disable client certificate authentication
  master_auth {
    client_certificate_config {
      issue_client_certificate = false
    }
  }

  # Enable Binary Authorization for container security
  binary_authorization {
    evaluation_mode = "PROJECT_SINGLETON_POLICY_ENFORCE"
  }

  # Configure Google Groups for RBAC management
  authenticator_groups_config {
    security_group = "gke-security-groups@${local.project_id}.iam.gserviceaccount.com"
  }

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
  location   = local.zone # Use zone instead of region
  node_count = var.min_nodes

  autoscaling {
    min_node_count = var.min_nodes
    max_node_count = var.max_nodes
  }

  node_config {
    machine_type = var.machine_type
    spot         = true          # Use spot instances for maximum cost savings
    disk_size_gb = 20            # Reduce disk size for cost savings
    disk_type    = "pd-standard" # Use standard disks instead of SSD

    metadata = {
      disable-legacy-endpoints = "true"
    }

    # Enable workload identity and GKE metadata server
    workload_metadata_config {
      mode = "GKE_METADATA"
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