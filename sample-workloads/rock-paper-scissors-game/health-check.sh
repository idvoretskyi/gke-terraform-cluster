#!/bin/bash

# Health Check Script for Rock Paper Scissors Game
# This script checks the health and performance of the deployed application

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE=${NAMESPACE:-default}
SERVICE_NAME="rock-paper-scissors-service"
APP_LABEL="app=rock-paper-scissors-game"

# Functions
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

# Check if kubectl is available
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed or not in PATH"
        exit 1
    fi
    
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Unable to connect to Kubernetes cluster"
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

# Check deployment status
check_deployment() {
    log_info "Checking deployment status..."
    
    local deployment_status=$(kubectl get deployment rock-paper-scissors-game -n $NAMESPACE -o jsonpath='{.status.conditions[?(@.type=="Available")].status}' 2>/dev/null || echo "NotFound")
    
    if [ "$deployment_status" = "True" ]; then
        log_success "Deployment is available"
    else
        log_error "Deployment is not available"
        kubectl get deployment rock-paper-scissors-game -n $NAMESPACE -o wide
        return 1
    fi
}

# Check pod status
check_pods() {
    log_info "Checking pod status..."
    
    local ready_pods=$(kubectl get pods -l $APP_LABEL -n $NAMESPACE -o jsonpath='{.items[?(@.status.phase=="Running")].metadata.name}' | wc -w)
    local total_pods=$(kubectl get pods -l $APP_LABEL -n $NAMESPACE --no-headers | wc -l)
    
    if [ $ready_pods -gt 0 ]; then
        log_success "$ready_pods/$total_pods pods are running"
    else
        log_error "No pods are running"
        kubectl get pods -l $APP_LABEL -n $NAMESPACE -o wide
        return 1
    fi
}

# Check service and external IP
check_service() {
    log_info "Checking service status..."
    
    local external_ip=$(kubectl get service $SERVICE_NAME -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
    
    if [ -n "$external_ip" ]; then
        log_success "Service has external IP: $external_ip"
        echo "Access the game at: http://$external_ip"
    else
        log_warning "Service external IP is not yet assigned"
        kubectl get service $SERVICE_NAME -n $NAMESPACE -o wide
    fi
}

# Check HPA status
check_hpa() {
    log_info "Checking HPA status..."
    
    local hpa_status=$(kubectl get hpa rock-paper-scissors-hpa -n $NAMESPACE -o jsonpath='{.status.conditions[?(@.type=="AbleToScale")].status}' 2>/dev/null || echo "NotFound")
    
    if [ "$hpa_status" = "True" ]; then
        log_success "HPA is functioning correctly"
        kubectl get hpa rock-paper-scissors-hpa -n $NAMESPACE
    else
        log_warning "HPA may not be configured correctly"
        kubectl get hpa rock-paper-scissors-hpa -n $NAMESPACE 2>/dev/null || log_info "HPA not found"
    fi
}

# Test application health endpoint
test_health_endpoint() {
    log_info "Testing application health endpoint..."
    
    local external_ip=$(kubectl get service $SERVICE_NAME -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
    
    if [ -n "$external_ip" ]; then
        local health_status=$(curl -s -o /dev/null -w "%{http_code}" http://$external_ip/health --max-time 10 || echo "000")
        
        if [ "$health_status" = "200" ]; then
            log_success "Health endpoint is responding correctly"
        else
            log_error "Health endpoint returned status: $health_status"
            return 1
        fi
    else
        log_warning "Cannot test health endpoint - no external IP available"
    fi
}

# Show resource usage
show_resource_usage() {
    log_info "Showing resource usage..."
    
    echo "Pod resource usage:"
    kubectl top pods -l $APP_LABEL -n $NAMESPACE 2>/dev/null || log_warning "Metrics not available (metrics-server may not be installed)"
    
    echo ""
    echo "Node resource usage:"
    kubectl top nodes 2>/dev/null || log_warning "Node metrics not available"
}

# Show recent logs
show_recent_logs() {
    log_info "Showing recent application logs..."
    
    kubectl logs -l $APP_LABEL -n $NAMESPACE --tail=10 --timestamps=true 2>/dev/null || log_warning "Unable to fetch logs"
}

# Main health check function
main() {
    echo "=== Rock Paper Scissors Game Health Check ==="
    echo "Namespace: $NAMESPACE"
    echo "Service: $SERVICE_NAME"
    echo "=========================================="
    echo ""
    
    local errors=0
    
    check_prerequisites || ((errors++))
    check_deployment || ((errors++))
    check_pods || ((errors++))
    check_service
    check_hpa
    test_health_endpoint || ((errors++))
    
    echo ""
    echo "=== Additional Information ==="
    show_resource_usage
    echo ""
    show_recent_logs
    
    echo ""
    echo "=========================================="
    if [ $errors -eq 0 ]; then
        log_success "All health checks passed!"
        exit 0
    else
        log_error "$errors health check(s) failed"
        exit 1
    fi
}

# Run health check
main "$@"
