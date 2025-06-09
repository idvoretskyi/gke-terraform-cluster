data "external" "username" {
  program = ["sh", "-c", "echo '{\"username\":\"'$(whoami)'\"}'"]
}

data "external" "gcloud_project" {
  program = ["sh", "-c", "PROJECT=$(gcloud config get-value project 2>/dev/null); if [ -z \"$PROJECT\" ]; then echo '{\"project_id\":\"\"}'; else echo '{\"project_id\":\"'$PROJECT'\"}'; fi"]
}

data "external" "gcloud_region" {
  program = ["sh", "-c", "REGION=$(gcloud config get-value compute/region 2>/dev/null); if [ -z \"$REGION\" ]; then echo '{\"region\":\"us-central1\"}'; else echo '{\"region\":\"'$REGION'\"}'; fi"]
}