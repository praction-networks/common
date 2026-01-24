#!/bin/bash

# Script to bump version for common package
# Usage: ./scripts/bump-version.sh [patch|minor|major]

set -e

# Fetch latest tags from remote
git fetch --tags --quiet

PACKAGE_NAME="github.com/praction-networks/common"

# Find the HIGHEST version tag from ALL tags (not just reachable from HEAD)
# This prevents collisions with tags created by other developers/CI
CURRENT_VERSION=$(git tag -l "v*" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -n 1)

if [[ -z "$CURRENT_VERSION" ]]; then
    CURRENT_VERSION="v0.0.0"
    echo "⚠️  No existing version tags found, starting from v0.0.0"
fi

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

# Check if tag exists
if git rev-parse "$NEXT_VERSION" >/dev/null 2>&1; then
    echo "❌ Error: Tag $NEXT_VERSION already exists locally or remotely (after fetch)."
    echo "Please pull latest changes or manually tag the correct version."
    exit 1
fi

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
