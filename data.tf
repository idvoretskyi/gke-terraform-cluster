# Get current username for cluster naming
data "external" "username" {
  program = ["sh", "-c", "echo '{\"username\":\"'$(whoami 2>/dev/null || echo 'default-user')'\"}'"]
}

# Get current gcloud project if available
data "external" "gcloud_project" {
  program = ["sh", "-c", "if command -v gcloud >/dev/null 2>&1; then PROJECT=$(gcloud config get-value project 2>/dev/null || echo ''); if [ -n \"$PROJECT\" ]; then echo '{\"project_id\":\"'$PROJECT'\"}'; else echo '{\"project_id\":\"\"}'; fi; else echo '{\"project_id\":\"\"}'; fi"]
}

# Get current gcloud region if available (fallback only)
data "external" "gcloud_region" {
  program = ["sh", "-c", "if command -v gcloud >/dev/null 2>&1; then REGION=$(gcloud config get-value compute/region 2>/dev/null || echo ''); if [ -n \"$REGION\" ]; then echo '{\"region\":\"'$REGION'\"}'; else echo '{\"region\":\"\"}'; fi; else echo '{\"region\":\"\"}'; fi"]
}