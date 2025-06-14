#!/bin/bash

# Terraform validation script
# This script ensures the configuration is valid before applying

set -e

echo "🔍 Checking Terraform configuration..."

# Initialize Terraform
echo "📦 Initializing Terraform..."
terraform init -backend=false

# Format check
echo "🎨 Checking format..."
terraform fmt -check=true

# Validate configuration
echo "✅ Validating configuration..."
terraform validate

# Run plan (dry-run)
echo "📋 Running plan (dry-run)..."
if terraform plan -out=tfplan; then
    echo "✅ Plan completed successfully"
    rm -f tfplan
else
    echo "❌ Plan failed"
    exit 1
fi

echo "🎉 All validations passed!"