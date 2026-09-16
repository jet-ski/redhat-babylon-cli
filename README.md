# babylon — CLI for Red Hat Demo Platform

A command-line interface for [demo.redhat.com](https://demo.redhat.com) (the Babylon platform) that provides the core features of the web UI for managing demo and lab services.

## Origin

This CLI was originally developed by [Guillaume Coré](https://github.com/fridim) as part of the [Babylon project](https://github.com/redhat-cop/babylon) in [PR #3210](https://github.com/redhat-cop/babylon/pull/3210). This repository extracts it as a standalone project to allow independent iteration while the upstream PR is pending.

## Commands

```
babylon
├── login                        # Browser-based OAuth login
├── catalog
│   ├── list                     # List available catalog items
│   └── describe <name>         # Show details of a catalog item
├── service (alias: svc)
│   ├── list                     # List your provisioned services
│   ├── status <name>           # Show service status
│   ├── order <item>            # Order a new service
│   ├── start <name>            # Start a stopped service
│   ├── stop <name>             # Stop a running service
│   ├── retire <name>           # Retire a service
│   ├── delete <name>           # Delete a service
│   └── delete-all [--yes]      # Delete all services
├── workshop
│   ├── list                     # List workshops
│   ├── create                   # Create a workshop
│   ├── delete                   # Delete a workshop
│   └── status                  # Show workshop status
└── version
```

## Install

### Build from source

```bash
make babylon
# Binary is in build/babylon
```

### Cross-compile

```bash
make babylon-cross
# Binaries for linux/darwin/windows in build/
```

## Usage

```bash
# Login (opens browser)
babylon login demo.redhat.com

# List available demos
babylon catalog list

# Order a service
babylon service order my-demo --end-date 72h --param region=us-east-1

# Check status
babylon service status my-demo-abc12

# Stop when done
babylon service stop my-demo-abc12
```

## Global Flags

| Flag | Env Var | Description | Default |
|------|---------|-------------|---------|
| `--server` | `BABYLON_SERVER` | Catalog API base URL | (from config) |
| `--namespace`, `-n` | `BABYLON_NAMESPACE` | Target namespace | User namespace |
| `--output`, `-o` | — | Output format: `table`, `json`, `yaml` | `table` |
| `--config` | — | Config file path | `~/.config/babylon/config.yaml` |

## Architecture

```
┌──────────────┐         ┌──────────────────┐         ┌──────────────┐
│              │  HTTPS   │                  │  K8s    │              │
│  babylon CLI ├─────────►│   Catalog API    ├────────►│  Kubernetes  │
│              │  REST    │  (Python/aiohttp)│  API    │  API Server  │
└──────────────┘         └──────────────────┘         └──────────────┘
```

The CLI talks exclusively to the Catalog API (the same backend the web UI uses). No `client-go` dependency, no direct Kubernetes access.

## Backend Requirement

The CLI login flow requires the `/auth/cli-redirect` endpoint on the Catalog API server. The patch for this endpoint is included in `patches/auth-cli-redirect.patch` for reference.

## License

Apache License 2.0 — same as the upstream Babylon project.

## Credits

- [Guillaume Coré (fridim)](https://github.com/fridim) — original CLI author
- [Babylon Project](https://github.com/redhat-cop/babylon) — upstream platform
