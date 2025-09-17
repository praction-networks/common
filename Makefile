# Makefile for common package

PACKAGE_NAME := github.com/praction-networks/common
CURRENT_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
SCRIPT_DIR := ./scripts

.PHONY: help version patch minor major tag push clean test build

## 📋 Show this help message
help:
	@echo "Common Package Management"
	@echo "========================"
	@echo ""
	@echo "Current version: $(CURRENT_VERSION)"
	@echo ""
	@echo "Available commands:"
	@echo "  version     - Show current version"
	@echo "  patch       - Bump patch version (v1.2.3 -> v1.2.4) + commit + push"
	@echo "  minor       - Bump minor version (v1.2.3 -> v1.3.0) + commit + push"
	@echo "  major       - Bump major version (v1.2.3 -> v2.0.0) + commit + push"
	@echo "  push        - Push changes and tags to remote"
	@echo "  test        - Run tests"
	@echo "  build       - Build the package"
	@echo "  clean       - Clean up generated files"
	@echo "  release-patch - Complete patch release (test + build + bump + commit + push + update deps)"
	@echo "  release-minor - Complete minor release (test + build + bump + commit + push + update deps)"
	@echo "  release-major - Complete major release (test + build + bump + commit + push + update deps)"

## 📊 Show current version
version:
	@echo "Current version: $(CURRENT_VERSION)"
	@echo "Use './scripts/bump-version.sh [patch|minor|major]' to bump version"

## 🔧 Bump patch version
patch:
	@$(SCRIPT_DIR)/bump-version.sh patch

## 🔧 Bump minor version
minor:
	@$(SCRIPT_DIR)/bump-version.sh minor

## 🔧 Bump major version
major:
	@$(SCRIPT_DIR)/bump-version.sh major

## 🏷️  Create and push git tag
tag:
	@echo "🏷️  Creating git tag for $(CURRENT_VERSION)..."
	@git tag -a $(CURRENT_VERSION) -m "Release $(CURRENT_VERSION)" 2>/dev/null || echo "Tag $(CURRENT_VERSION) already exists"
	@echo "✅ Tag created: $(CURRENT_VERSION)"

## 📤 Push changes and tags to remote
push:
	@echo "📤 Pushing changes and tags to remote..."
	@git push origin main
	@git push origin --tags
	@echo "✅ Pushed to remote"

## 🧪 Run tests
test:
	@echo "🧪 Running tests..."
	@go test ./...

## 🏗️  Build the package
build:
	@echo "🏗️  Building package..."
	@go mod tidy
	@go mod verify
	@echo "✅ Package built successfully"

## 🧹 Clean up generated files
clean:
	@echo "🧹 Cleaning up..."
	@go clean -cache
	@echo "✅ Cleaned up"

## 🚀 Complete patch release
release-patch: test build
	@$(SCRIPT_DIR)/bump-version.sh patch
	@make update-deps
	@echo "🎉 Patch release completed!"

## 🚀 Complete minor release
release-minor: test build
	@$(SCRIPT_DIR)/bump-version.sh minor
	@make update-deps
	@echo "🎉 Minor release completed!"

## 🚀 Complete major release
release-major: test build
	@$(SCRIPT_DIR)/bump-version.sh major
	@make update-deps
	@echo "🎉 Major release completed!"

## 📦 Update all dependent services to use new version
update-deps:
	@$(SCRIPT_DIR)/update-deps.sh
