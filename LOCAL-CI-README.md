# Local CI Development Guide

Run your GitLab CI pipeline locally before pushing — no more waiting for remote pipelines!

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Available Jobs](#available-jobs)
- [How It Works](#how-it-works)
- [Private Bitbucket Dependencies](#private-bitbucket-dependencies)
- [Enabling Full Pipeline (SSH Keys)](#enabling-full-pipeline-ssh-keys)
- [Files Overview](#files-overview)
- [Tested Output](#tested-output)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

| Tool | Required | Purpose |
|------|----------|---------|
| **Docker** | Yes | Runs CI jobs in containers |
| **gitlab-ci-local** | Yes | Executes GitLab CI pipelines locally |
| **Git** | Yes | Version control |

> **Note**: Windows users must use **WSL2** (Windows Subsystem for Linux).

---

## Installation

### Option A: Ubuntu / Debian / EC2

```bash
# 1. Update system
sudo apt update && sudo apt upgrade -y

# 2. Install Docker
sudo apt install -y docker.io
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
newgrp docker

# 3. Verify Docker
docker --version

# 4. Install gitlab-ci-local (via npm)
sudo apt install -y nodejs npm
sudo npm install -g gitlab-ci-local

# 5. Verify gitlab-ci-local
gitlab-ci-local --version
```

### Option B: macOS

```bash
# 1. Install Docker Desktop from https://www.docker.com/products/docker-desktop/
# 2. Install gitlab-ci-local
brew install gitlab-ci-local
```

### Option C: Windows (via WSL2)

```powershell
# 1. Open PowerShell as Admin
wsl --install

# 2. Inside WSL, follow Ubuntu instructions above
```

---

## Quick Start

```bash
# 1. Clone the repository
git clone <your-repo-url>
cd go-main

# 2. List all available CI jobs
gitlab-ci-local --file .gitlab-ci-local.yml --list

# 3. Run Go Lint (always works, no credentials needed)
gitlab-ci-local --file .gitlab-ci-local.yml "Go Lint"

# 4. Run all jobs at once
gitlab-ci-local --file .gitlab-ci-local.yml
```

---

## Available Jobs

| Job | Stage | Requires Bitbucket SSH Key? | Expected Result (No Creds) | Expected Result (With Creds) |
|-----|-------|-----------------------------|---------------------------|------------------------------|
| **Go Lint** | lint | No | PASS | PASS |
| **Go Unit Test** | test | Yes | FAIL (expected) | PASS |
| **Go Vet** | test | Yes | FAIL (expected) | PASS |
| **Go Vulncheck** | test | Yes | FAIL (expected) | PASS |

### Run Individual Jobs

```bash
# Go Lint — formatting, spelling, line length checks
gitlab-ci-local --file .gitlab-ci-local.yml "Go Lint"

# Go Unit Test — runs all unit tests with coverage
gitlab-ci-local --file .gitlab-ci-local.yml "Go Unit Test"

# Go Vet — reports suspicious constructs
gitlab-ci-local --file .gitlab-ci-local.yml "Go Vet"

# Go Vulncheck — scans for known vulnerabilities
gitlab-ci-local --file .gitlab-ci-local.yml "Go Vulncheck"
```

---

## How It Works

### The Problem

The real `.gitlab-ci.yml` uses **remote GitLab Components**:

```yaml
include:
  - component: gitlab.com/d5100/.../go@$CI_COMMIT_SHA
```

These components reference `$CI_SERVER_FQDN` and `$CI_PROJECT_PATH`, which do **not resolve** outside GitLab. This means developers cannot test CI pipelines locally before pushing.

### The Solution

`.gitlab-ci-local.yml` is a **local equivalent** that:

1. **Replaces remote components** with local job definitions
2. **Uses public Docker images** (`golang:1.24`) instead of private AWS ECR images
3. **Mirrors the same job names** so developers get the same feedback locally
4. **Go Lint runs without any dependencies** — uses `gofmt`, `misspell`, and line-length checks directly (no `golangci-lint` which requires package loading)

### What Go Lint Checks Locally

| Check | Tool | What It Catches |
|-------|------|-----------------|
| Code formatting | `gofmt -l -s` | Improperly formatted Go code |
| Spelling errors | `misspell` | Typos in code and comments |
| Line length | `awk` (140 char max) | Lines exceeding 140 characters |

> **Why not golangci-lint?** — golangci-lint v2 **always** loads Go packages internally (even for simple linters). Since `go.mod` references private `bitbucket.org/csgot/*` packages, package loading fails without credentials. Direct tools bypass this entirely.

---

## Private Bitbucket Dependencies

The `go.mod` file requires private packages:

```
bitbucket.org/csgot/helis-elasticlogrus
bitbucket.org/csgot/helis-go-uuid
bitbucket.org/csgot/helis-market-base-client-go
bitbucket.org/csgot/helis-market-settings-client-go
```

These packages need **Bitbucket SSH access** to download. Without credentials:
- **Go Lint** works (does not download packages)
- **Go Unit Test / Vet / Vulncheck** fail at `go mod download` (set to `allow_failure: true`)

---

## Enabling Full Pipeline (SSH Keys)

Once you have SSH access to `bitbucket.org/csgot`:

```bash
# 1. Generate SSH key (if you don't have one)
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa -N ""

# 2. Copy your public key
cat ~/.ssh/id_rsa.pub

# 3. Add it to your Bitbucket account:
#    Bitbucket → Settings → SSH Keys → Add Key

# 4. Test SSH connection
ssh -T git@bitbucket.org

# 5. Uncomment GIT_SSH_KEY in .gitlab-ci-local-variables.yml:
#    GIT_SSH_KEY: "~/.ssh/id_rsa"

# 6. Run all jobs — they should all pass now
gitlab-ci-local --file .gitlab-ci-local.yml
```

---

## Files Overview

```
go-main/
├── .gitlab-ci.yml                  ← Real GitLab CI (uses remote components)
├── .gitlab-ci-local.yml            ← Local CI wrapper (run with gitlab-ci-local)
├── .gitlab-ci-local-variables.yml  ← Local variable overrides
├── LOCAL-CI-README.md              ← This documentation file
└── tests/app/
    ├── .golangci.yml               ← Full linter config (used on GitLab CI)
    └── .golangci-local.yml         ← Lightweight linter config (for local use)
```

| File | Purpose |
|------|---------|
| `.gitlab-ci-local.yml` | Main file — pass this to `gitlab-ci-local --file` |
| `.gitlab-ci-local-variables.yml` | Variable overrides (SSH key config goes here) |
| `.golangci.yml` | Production linter config (30+ linters, runs on GitLab) |
| `.golangci-local.yml` | Local linter config (syntax-only, no package loading) |

---

## Tested Output

Below is the actual output from a successful run on an EC2 Ubuntu instance:

```
$ gitlab-ci-local --file .gitlab-ci-local.yml "Go Lint"

Go Lint starting golang:1.24 (lint)
Go Lint $ cd tests/app
Go Lint $ echo "=== Running gofmt ==="
Go Lint > === Running gofmt ===
Go Lint > PASS: All files properly formatted
Go Lint $ echo "=== Installing misspell ==="
Go Lint > === Installing misspell ===
Go Lint $ go install github.com/client9/misspell/cmd/misspell@latest
Go Lint $ echo "=== Running misspell ==="
Go Lint > === Running misspell ===
Go Lint > PASS: No misspellings found
Go Lint $ echo "=== Checking line length (max 140 chars) ==="
Go Lint > === Checking line length (max 140 chars) ===
Go Lint > PASS: All lines within 140 characters
Go Lint $ echo "=== Local lint checks complete ==="
Go Lint > === Local lint checks complete ===
Go Lint > ALL CHECKS PASSED
Go Lint finished in 48 s

 PASS  Go Lint
```

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| `Cannot connect to Docker` | Start Docker: `sudo systemctl start docker` |
| `go mod download` fails | Expected without Bitbucket SSH — Go Lint still works |
| `permission denied` on Docker | Run: `sudo usermod -aG docker $USER && newgrp docker` |
| `gitlab-ci-local: command not found` | Install: `sudo npm install -g gitlab-ci-local` |
| Jobs take too long first time | Docker images download on first run — subsequent runs are faster |
| Predefined vars warning | Safe to ignore — does not affect job execution |
