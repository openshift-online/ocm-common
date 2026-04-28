# ocm-common
This project contains utility functions intended for sharing across OCM clients and services. Its objective is to streamline OCM codebases by eliminating duplicate code found within them.

## Quick Start

```bash
make build       # Build the project
make test        # Run all tests
make lint        # Run golangci-lint
```

## How to run the tests
* `make test` -- runs unit tests

## Contributing
[Contribution guide](CONTRIBUTING.md)

## Installation

```bash
go get github.com/openshift-online/ocm-common
```

Import the packages you need:

```go
import "github.com/openshift-online/ocm-common/pkg/aws"
import "github.com/openshift-online/ocm-common/pkg/utils"
```

## Usage

This library provides utility functions organized by domain:

- `pkg/aws/` — AWS SDK helpers (STS, EC2, IAM operations)
- `pkg/gcp/` — GCP SDK helpers
- `pkg/cluster/` — Cluster-related utilities
- `pkg/utils/` — General-purpose helpers

Refer to the [GoDoc](https://pkg.go.dev/github.com/openshift-online/ocm-common) for detailed API documentation.
