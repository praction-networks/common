# Making Common Package Private

## ✅ Steps to Make the Package Private

### 1. **Set Go Environment Variables**

Add to your `~/.bashrc` or `~/.zshrc`:
```bash
export GOPRIVATE=github.com/praction-networks/*
export GONOPROXY=github.com/praction-networks/*
export GONOSUMDB=github.com/praction-networks/*
```

Or source the project `.envrc` file:
```bash
source /home/praction/development/hotspot/.envrc
```

### 2. **Configure Git Authentication**

#### Option A: SSH (Recommended)
```bash
# Configure Git to use SSH for private repos
git config --global url."git@github.com:praction-networks/".insteadOf "https://github.com/praction-networks/"

# Ensure SSH key is added to GitHub
ssh -T git@github.com
```

#### Option B: Personal Access Token
```bash
# Set up Git credential helper
git config --global credential.helper store

# Or use token in URL (less secure)
git config --global url."https://YOUR_TOKEN@github.com/praction-networks/".insteadOf "https://github.com/praction-networks/"
```

### 3. **Make GitHub Repository Private**

1. Go to: https://github.com/praction-networks/common/settings
2. Scroll to "Danger Zone"
3. Click "Change visibility" → "Make private"

### 4. **Verify Configuration**

```bash
# Check Go environment
go env GOPRIVATE GONOPROXY GONOSUMDB

# Test fetching the package
go get github.com/praction-networks/common@latest
```

### 5. **Update All Services**

After making the repo private, update all services:
```bash
cd common
./scripts/update-deps.sh
```

## 🔒 Security Benefits

- ✅ Prevents code from being fetched from public proxies
- ✅ Skips checksum verification (private repos don't need it)
- ✅ Ensures only authenticated users can access
- ✅ Protects proprietary code

## 📝 Notes

- The `replace` directive in `go.mod` still works for local development
- All services using the common package need these environment variables
- CI/CD pipelines also need these variables configured

