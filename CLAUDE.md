# CLAUDE.md

<!-- Canonical source: AGENTS.md. This file is auto-generated for Claude Code compatibility. -->

This file provides guidance to AI coding assistants when working with this repository.

## Project Overview

OCM Common — a shared Go library containing utility functions used across OCM clients and services. Provides reusable helpers for AWS, GCP, cloud provider interactions, networking, and general-purpose utilities.

## Build & Test Commands

```bash
make build           # Build the project
make test            # Run all tests
make lint            # Run golangci-lint
make coverage        # Generate test coverage report
make fmt             # Format Go source code
make clean           # Remove build artifacts
```

## Architecture

- **pkg/**: All library code organized by domain
  - **pkg/aws/**: AWS SDK helpers and utilities
  - **pkg/gcp/**: GCP SDK helpers and utilities
  - **pkg/cluster/**: Cluster-related utilities
  - **pkg/utils/**: General-purpose utility functions

## Key Conventions

- Module path: `github.com/openshift-online/ocm-common`
- Pure library — no main package
- All public APIs should be well-documented with Go doc comments
- Maintain backward compatibility for exported symbols
