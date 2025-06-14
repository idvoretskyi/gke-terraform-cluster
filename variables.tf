variable "project_id" {
  description = "The GCP project ID where the cluster will be created (defaults to current gcloud project)"
  type        = string
  default     = null
}

variable "region" {
  description = "The GCP region where the cluster will be created (defaults to current gcloud region)"
  type        = string
  default     = null
}

variable "cluster_name_suffix" {
  description = "Optional suffix for the cluster name (username prefix will be automatically added)"
  type        = string
  default     = "gke-cluster"
}

variable "machine_type" {
  description = "The machine type for the GKE nodes"
  type        = string
  default     = "e2-micro"
}

variable "min_nodes" {
  description = "The minimum number of nodes in the node pool"
  type        = number
  default     = 1
}

variable "max_nodes" {
  description = "The maximum number of nodes in the node pool"
  type        = number
  default     = 3
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"

  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod."
  }
}

variable "enable_private_cluster" {
  description = "Enable private cluster configuration"
  type        = bool
  default     = true
}

variable "enable_private_endpoint" {
  description = "Enable private endpoint (requires bastion host or VPN)"
  type        = bool
  default     = false # Changed to false for easier development access
}

variable "master_authorized_networks" {
  description = "List of master authorized networks"
  type = list(object({
    cidr_block   = string
    display_name = string
  }))
  default = [
    {
      cidr_block   = "0.0.0.0/0"
      display_name = "All networks - development only"
    }
  ]
}

variable "enable_network_policy" {
  description = "Enable network policy for enhanced security"
  type        = bool
  default     = true
}

variable "enable_binary_authorization" {
  description = "Enable Binary Authorization for container security"
  type        = bool
  default     = false # Changed to false for easier development
}

variable "release_channel" {
  description = "GKE release channel (RAPID, REGULAR, STABLE)"
  type        = string
  default     = "REGULAR"

  validation {
    condition     = contains(["RAPID", "REGULAR", "STABLE"], var.release_channel)
    error_message = "Release channel must be one of: RAPID, REGULAR, STABLE."
  }
}

variable "node_pool_disk_size" {
  description = "Disk size for node pool (GB)"
  type        = number
  default     = 20
}

variable "node_pool_disk_type" {
  description = "Disk type for node pool (pd-standard, pd-balanced, pd-ssd)"
  type        = string
  default     = "pd-standard"
}

variable "enable_spot_instances" {
  description = "Enable spot instances for cost optimization"
  type        = bool
  default     = true
}

variable "resource_labels" {
  description = "Resource labels for cost tracking and management"
  type        = map(string)
  default = {
    managed-by = "terraform"
  }
}