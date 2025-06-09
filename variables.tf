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