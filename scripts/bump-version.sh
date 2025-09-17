#!/bin/bash

# Script to bump version for common package
# Usage: ./scripts/bump-version.sh [patch|minor|major]

set -e

PACKAGE_NAME="github.com/praction-networks/common"
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Extract version numbers
VERSION_REGEX="v([0-9]+)\.([0-9]+)\.([0-9]+)"
if [[ $CURRENT_VERSION =~ $VERSION_REGEX ]]; then
    MAJOR=${BASH_REMATCH[1]}
    MINOR=${BASH_REMATCH[2]}
    PATCH=${BASH_REMATCH[3]}
else
    echo "❌ Could not parse current version: $CURRENT_VERSION"
    exit 1
fi

# Calculate next version based on argument
case "${1:-patch}" in
    "patch")
        NEXT_VERSION="v${MAJOR}.${MINOR}.$((PATCH + 1))"
        ;;
    "minor")
        NEXT_VERSION="v${MAJOR}.$((MINOR + 1)).0"
        ;;
    "major")
        NEXT_VERSION="v$((MAJOR + 1)).0.0"
        ;;
    *)
        echo "❌ Invalid version type: $1"
        echo "Usage: $0 [patch|minor|major]"
        exit 1
        ;;
esac

echo "📊 Current version: $CURRENT_VERSION"
echo "🚀 Next version: $NEXT_VERSION"

# Check if there are any changes to commit
if [[ -n $(git status --porcelain) ]]; then
    echo "📝 Found uncommitted changes:"
    git status --short
    echo ""
    
    # Add all changes
    echo "📦 Adding all changes..."
    git add .
    
    # Commit changes
    echo "💾 Committing changes..."
    git commit -m "chore: prepare for release $NEXT_VERSION"
else
    echo "ℹ️  No changes to commit"
fi

# Confirm before proceeding
read -p "Do you want to create tag $NEXT_VERSION and push to remote? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Cancelled"
    exit 1
fi

# Create the tag
echo "🏷️  Creating tag $NEXT_VERSION..."
git tag -a "$NEXT_VERSION" -m "Release $NEXT_VERSION"

# Push changes and tags
echo "📤 Pushing changes and tags to remote..."
git push origin main
git push origin --tags

echo "✅ Release $NEXT_VERSION completed successfully!"
echo ""
echo "Next steps:"
echo "1. Run 'make update-deps' to update all dependent services"
