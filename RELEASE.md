# Common Package Release Management

This document explains how to manage versions and releases for the common package.

## Quick Start

```bash
# Show current version
make version

# Bump patch version (v1.2.3 -> v1.2.4) + commit + push
make patch

# Bump minor version (v1.2.3 -> v1.3.0) + commit + push
make minor

# Bump major version (v1.2.3 -> v2.0.0) + commit + push
make major

# Complete release (test + build + bump + commit + push + update deps)
make release-patch
make release-minor
make release-major
```

## Available Commands

### Version Management
- `make version` - Show current version information
- `make patch` - Bump patch version + commit + push (backward-compatible bug fixes)
- `make minor` - Bump minor version + commit + push (backward-compatible new features)
- `make major` - Bump major version + commit + push (breaking changes)

### Release Process
- `make release-patch` - Complete patch release (test + build + bump + commit + push + update deps)
- `make release-minor` - Complete minor release (test + build + bump + commit + push + update deps)
- `make release-major` - Complete major release (test + build + bump + commit + push + update deps)

### Development
- `make test` - Run tests
- `make build` - Build and verify the package
- `make clean` - Clean up generated files

### Publishing
- `make push` - Push changes and tags to remote repository
- `make update-deps` - Update all dependent services to use the new version

## Release Workflow

### 1. Patch Release (Bug Fixes)
```bash
# For backward-compatible bug fixes
make release-patch
```

### 2. Minor Release (New Features)
```bash
# For backward-compatible new features
make release-minor
```

### 3. Major Release (Breaking Changes)
```bash
# For breaking changes
make release-major
```

### 4. Simple Version Bump (without full release)
If you just want to bump version and push without running tests:
```bash
make patch    # or minor/major
```

## Manual Process

If you prefer to do it manually:

1. **Test and Build**
   ```bash
   make test
   make build
   ```

2. **Bump Version**
   ```bash
   # Using the script directly
   ./scripts/bump-version.sh patch
   # or
   ./scripts/bump-version.sh minor
   # or
   ./scripts/bump-version.sh major
   ```

3. **Push Changes**
   ```bash
   make push
   ```

4. **Update Dependencies**
   ```bash
   make update-deps
   ```

## Version Numbering

We follow [Semantic Versioning](https://semver.org/):
- **MAJOR** version when you make incompatible API changes
- **MINOR** version when you add functionality in a backward-compatible manner
- **PATCH** version when you make backward-compatible bug fixes

## Git Tags

Each release creates a Git tag in the format `vX.Y.Z` (e.g., `v1.2.3`). These tags are used by Go modules to reference specific versions of the package.

## Dependencies

The common package is used by multiple services. After releasing a new version, run `make update-deps` to update all dependent services to use the latest version.
