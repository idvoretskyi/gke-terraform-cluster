variable "project_id" {
  description = "The GCP project ID where the cluster will be created"
  type        = string
}

variable "region" {
  description = "The GCP region where the cluster will be created"
  type        = string
  default     = "europe-west4"
}

variable "cluster_name" {
  description = "The base name of the GKE cluster"
  type        = string
  default     = "cloud-phoenix"  # You can change this to any fictional name you like
}

variable "machine_type" {
  description = "The machine type for the GKE nodes"
  type        = string
  default     = "e2-medium"
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