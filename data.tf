data "external" "username" {
  program = ["sh", "-c", "echo '{\"username\":\"'$(whoami 2>/dev/null || echo 'github-actions')'\"}'"]
}

data "external" "gcloud_project" {
  program = ["sh", "-c", "if command -v gcloud >/dev/null 2>&1; then PROJECT=$(gcloud config get-value project 2>/dev/null); if [ -z \"$PROJECT\" ]; then echo '{\"project_id\":\"\"}'; else echo '{\"project_id\":\"'$PROJECT'\"}'; fi; else echo '{\"project_id\":\"\"}'; fi"]
}

data "external" "gcloud_region" {
  program = ["sh", "-c", "if command -v gcloud >/dev/null 2>&1; then REGION=$(gcloud config get-value compute/region 2>/dev/null); if [ -z \"$REGION\" ]; then echo '{\"region\":\"\"}'; else echo '{\"region\":\"'$REGION'\"}'; fi; else echo '{\"region\":\"\"}'; fi"]
}