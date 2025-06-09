# GKE Terraform Cluster

## Description

This Terraform module creates a Google Kubernetes Engine (GKE) cluster with autoscaling capabilities, preemptible nodes, and proper security configurations.

## Features

- **Ultra Cost-Optimized**: Uses e2-micro instances, spot VMs, and minimal disk sizes
- **Preemptible + Spot Instances**: Maximum cost savings with preemptible and spot nodes
- **Minimal Resource Usage**: 20GB standard disks, reduced OAuth scopes
- **Dynamic Naming**: Cluster names use your username with random suffix
- **Auto-scaling**: Configurable min/max node counts for cost control
- **Security**: Best practices with disabled legacy endpoints

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
- GCP project with the following APIs enabled:
  - Kubernetes Engine API
  - Compute Engine API

## Quick Start

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd gke-terraform-cluster
   ```

2. **Authenticate with GCP:**
   ```bash
   gcloud auth application-default login
   ```

3. **Initialize Terraform:**
   ```bash
   terraform init
   ```

4. **Plan the deployment:**
   ```bash
   # Auto-detects current gcloud project
   terraform plan
   
   # Or specify a different project
   terraform plan -var="project_id=your-gcp-project-id"
   ```

5. **Apply the configuration:**
   ```bash
   # Auto-detects current gcloud project
   terraform apply
   
   # Or specify a different project
   terraform apply -var="project_id=your-gcp-project-id"
   ```

The module will automatically:
- **Auto-detect your current gcloud project** (can be overridden with project_id variable)
- Create a cluster name prefixed with your username (e.g., `idv-gke-cluster-abc123`)
- Use the detected or specified project_id for all GCP resources

## Configuration

### Variables

The following variables can be configured:

| Variable | Description | Type | Default |
|----------|-------------|------|---------|
| `project_id` | The GCP project ID where the cluster will be created (defaults to current gcloud project) | `string` | Auto-detected |
| `region` | The GCP region where the cluster will be created | `string` | `"us-central1"` |
| `cluster_name_suffix` | Suffix for the cluster name (username prefix added automatically) | `string` | `"gke-cluster"` |
| `machine_type` | The machine type for the GKE nodes (cost-optimized default) | `string` | `"e2-micro"` |
| `min_nodes` | The minimum number of nodes in the node pool | `number` | `1` |
| `max_nodes` | The maximum number of nodes in the node pool | `number` | `3` |

### Example terraform.tfvars

```hcl
# Optional: Override auto-detected project
# project_id           = "my-gcp-project-id"

# Optional variables - these are just examples of customizations
region               = "us-central1"
cluster_name_suffix  = "my-cluster"
machine_type         = "e2-medium"
min_nodes           = 2
max_nodes           = 5
```

## Outputs

| Output | Description |
|--------|-------------|
| `cluster_name` | The name of the GKE cluster |
| `cluster_endpoint` | The endpoint of the GKE cluster |
| `cluster_ca_certificate` | The cluster CA certificate (base64 encoded) |
| `cluster_location` | The location of the GKE cluster |
| `cluster_id` | The ID of the GKE cluster |
| `master_version` | The version of Kubernetes used by the GKE cluster |
| `node_pool_instance_group_urls` | List of instance group URLs for the node pools |
| `kubeconfig_command` | Command to configure kubectl for this cluster |
| `project_id` | The GCP project ID being used |

## Connecting to the Cluster

After deployment, use the output command to configure kubectl:

```bash
# Get the command from Terraform output
terraform output kubeconfig_command

# Or manually:
gcloud container clusters get-credentials <cluster-name> --region <region> --project <project-id>
```

## Cleanup

To destroy the cluster and all associated resources:

```bash
terraform destroy
```

## Cost Optimization Features

- **e2-micro instances**: Smallest available machine type for maximum savings
- **Preemptible + Spot VMs**: Up to 91% cost savings vs regular instances
- **20GB standard disks**: Minimal disk size with standard (not SSD) storage
- **Minimal OAuth scopes**: Only essential permissions to reduce attack surface
- **STABLE release channel**: Predictable updates, no premium for latest features
- **No cluster autoscaling**: Simplified scaling to avoid unexpected costs
- **Resource labels**: Cost tracking and management tags

## Security Considerations

- Legacy endpoints are disabled for security
- Auto-repair and auto-upgrade are enabled
- Minimal OAuth scopes for reduced attack surface
- Resource labels for compliance and tracking

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the [MIT License](LICENSE).
