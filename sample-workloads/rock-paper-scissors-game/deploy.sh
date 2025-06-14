#!/bin/bash

# Rock Paper Scissors Game Deployment Script
# This script builds and deploys the game to Google Kubernetes Engine

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
DEFAULT_VERSION="latest"
DEFAULT_PLATFORM="linux/amd64"

# Functions
print_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -p, --project PROJECT_ID    GCP Project ID (required)"
    echo "  -v, --version VERSION       Docker image version (default: $DEFAULT_VERSION)"
    echo "  -t, --tag VERSION           Alias for --version"
    echo "  --platform PLATFORM         Docker platform (default: $DEFAULT_PLATFORM)"
    echo "  --build-only                Only build and push, don't deploy"
    echo "  --deploy-only               Only deploy (assumes image exists)"
    echo "  --auto-version              Use git commit hash as version"
    echo "  -h, --help                  Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 -p my-gcp-project"
    echo "  $0 -p my-gcp-project -v v1.2.3"
    echo "  $0 -p my-gcp-project --auto-version"
    echo "  $0 -p my-gcp-project --build-only"
}

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Parse command line arguments
PROJECT_ID=""
VERSION="$DEFAULT_VERSION"
PLATFORM="$DEFAULT_PLATFORM"
BUILD_ONLY=false
DEPLOY_ONLY=false
AUTO_VERSION=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--project)
            PROJECT_ID="$2"
            shift 2
            ;;
        -v|--version|-t|--tag)
            VERSION="$2"
            shift 2
            ;;
        --platform)
            PLATFORM="$2"
            shift 2
            ;;
        --build-only)
            BUILD_ONLY=true
            shift
            ;;
        --deploy-only)
            DEPLOY_ONLY=true
            shift
            ;;
        --auto-version)
            AUTO_VERSION=true
            shift
            ;;
        -h|--help)
            print_usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

# Validate required parameters
if [[ -z "$PROJECT_ID" ]]; then
    log_error "Project ID is required. Use -p or --project flag."
    print_usage
    exit 1
fi

# Auto-generate version if requested
if [[ "$AUTO_VERSION" == true ]]; then
    if command -v git &> /dev/null && git rev-parse --git-dir > /dev/null 2>&1; then
        VERSION=$(git rev-parse --short HEAD)
        log_info "Using auto-generated version: $VERSION"
    else
        log_warning "Git not available or not in git repository. Using default version: $VERSION"
    fi
fi

# Set image name
IMAGE_NAME="gcr.io/$PROJECT_ID/rock-paper-scissors-game:$VERSION"

log_info "🎮 Rock Paper Scissors Game Deployment"
log_info "Project ID: $PROJECT_ID"
log_info "Version: $VERSION"
log_info "Image: $IMAGE_NAME"
log_info "Platform: $PLATFORM"

# Build and push Docker image
if [[ "$DEPLOY_ONLY" != true ]]; then
    log_info "🔨 Building and pushing Docker image..."
    
    if docker buildx build --platform "$PLATFORM" -t "$IMAGE_NAME" . --push; then
        log_success "Docker image built and pushed successfully"
    else
        log_error "Failed to build and push Docker image"
        exit 1
    fi
fi

# Deploy to Kubernetes
if [[ "$BUILD_ONLY" != true ]]; then
    log_info "🚀 Deploying to Kubernetes..."
    
    # Create a temporary deployment file with substituted values
    TEMP_DEPLOYMENT=$(mktemp)
    
    # Replace placeholders in deployment.yaml
    sed -e "s/PROJECT_ID/$PROJECT_ID/g" \
        -e "s/VERSION_TAG/$VERSION/g" \
        deployment.yaml > "$TEMP_DEPLOYMENT"
    
    if kubectl apply -f "$TEMP_DEPLOYMENT"; then
        log_success "Application deployed successfully"
        
        # Clean up temporary file
        rm "$TEMP_DEPLOYMENT"
        
        log_info "📋 Getting service information..."
        kubectl get svc rock-paper-scissors-service
        
        log_info "🔍 Getting pod status..."
        kubectl get pods -l app=rock-paper-scissors-game
        
        log_success "🎉 Deployment complete! Check the external IP above to access your game."
    else
        log_error "Failed to deploy application"
        rm "$TEMP_DEPLOYMENT"
        exit 1
    fi
fi

log_success "✅ All operations completed successfully!"
