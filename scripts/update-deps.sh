#!/bin/bash

# Script to update all dependent services to use the latest version of common package

set -e

PACKAGE_NAME="github.com/praction-networks/common"
LATEST_VERSION=$(git describe --tags --abbrev=0)

echo "📦 Updating dependencies in all services..."
echo "Latest version: $LATEST_VERSION"

# Find all go.mod files except the one in common directory
find .. -name "go.mod" -not -path "./go.mod" -not -path "../common/go.mod" | while read -r go_mod_file; do
    service_dir=$(dirname "$go_mod_file")
    service_name=$(basename "$service_dir")
    
    echo "Updating $service_name to $PACKAGE_NAME@$LATEST_VERSION..."
    
    cd "$service_dir"
    go get -u "$PACKAGE_NAME@$LATEST_VERSION"
    go mod tidy
    cd - > /dev/null
done

echo "✅ Dependencies updated in all services"
