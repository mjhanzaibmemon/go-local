# Local CI Development Guide

Run your GitLab CI pipeline locally before pushing — no more waiting for remote pipelines!

---

## Prerequisites

| Tool | Install |
|------|---------|
| **Docker Desktop** | [docker.com/products/docker-desktop](https://www.docker.com/products/docker-desktop/) |
| **gitlab-ci-local** | `brew install gitlab-ci-local` |

> **Note**: `gitlab-ci-local` requires macOS or Linux. Windows users should use WSL2.

---

## Quick Start

```bash
# 1. Go to repo root
cd go-main/

# 2. List all available jobs
gitlab-ci-local --file .gitlab-ci-local.yml --list

# 3. Run a specific job
gitlab-ci-local --file .gitlab-ci-local.yml "Go Lint"
gitlab-ci-local --file .gitlab-ci-local.yml "Go Unit Test"
gitlab-ci-local --file .gitlab-ci-local.yml "Go Vet"
gitlab-ci-local --file .gitlab-ci-local.yml "Go Vulncheck"

# 4. Run all jobs
gitlab-ci-local --file .gitlab-ci-local.yml
```

---

## Available Jobs

| Job | Stage | Needs Private Bitbucket Creds? | Expected Result |
|-----|-------|-------------------------------|-----------------|
| **Go Lint** | lint | ❌ No | ✅ Always passes |
| **Go Unit Test** | test | ✅ Yes | ⚠️ Fails without creds |
| **Go Vet** | test | ✅ Yes | ⚠️ Fails without creds |
| **Go Vulncheck** | test | ✅ Yes | ⚠️ Fails without creds |

---

## About Private Bitbucket Dependencies

The `go.mod` file requires private packages from `bitbucket.org/csgot/*`.

**Go Lint** does NOT need these — it will always work locally ✅

**Go Unit Test / Vet / Vulncheck** need credentials to run `go mod download`.

Once you have SSH access to `bitbucket.org/csgot`, add your key path to `.gitlab-ci-local-variables.yml`:

```yaml
GIT_SSH_KEY: "~/.ssh/id_rsa"
```

---

## How This Works

The real `.gitlab-ci.yml` uses remote GitLab Components that cannot be fetched locally:

```yaml
include:
  - component: gitlab.com/d5100/.../go@$CI_COMMIT_SHA   # ← not accessible locally
```

`.gitlab-ci-local.yml` replaces these with local equivalents:
- Defines the same job names (`Go Lint`, `Go Unit Test`, etc.)
- Uses public Docker images instead of private AWS ECR images
- Mirrors the same script logic from `templates/go.yml`

---

## Files in This Setup

```
go-main/
├── .gitlab-ci-local.yml            ← Run this file with gitlab-ci-local
├── .gitlab-ci-local-variables.yml  ← Local CI variable overrides
└── LOCAL-CI-README.md              ← This file
```

---

## Troubleshooting

**"Cannot connect to Docker"**
→ Start Docker Desktop and wait for it to fully load.

**"go mod download" fails**
→ You need Bitbucket SSH credentials for `bitbucket.org/csgot/*`. This is expected — Go Lint still works.

**"golangci-lint" lint errors**
→ Real lint issues found in your code. Fix them before pushing!
