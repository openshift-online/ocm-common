# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

OCM-Common is a Go library providing shared utility functions for OpenShift Cluster Manager (OCM) clients and services. It eliminates code duplication across OCM codebases by centralizing common functionality for AWS operations, cluster management, validation, and testing utilities.

## Development Commands

### Build and Test
- `make build` - Build the project
- `make test` - Run all unit tests
- `make coverage` - Run tests with coverage report
- `make fmt` - Format code using gofmt
- `make lint` - Run golangci-lint with 5-minute timeout
- `make clean` - Clean build artifacts

### Testing Framework
- Uses Ginkgo for all testing (third-party testing packages are rejected)
- Test files follow `*_test.go` naming convention
- Test suites use `*_suite_test.go` naming convention

## Code Architecture

### Core Package Structure

#### AWS Integration (`pkg/aws/`)
- **aws_client/**: Centralized AWS SDK v2 client with support for multiple AWS services (EC2, Route53, IAM, KMS, CloudFormation, etc.)
- **consts/**: AWS-related constants and default values
- **validations/**: AWS resource validation utilities including IAM helpers and tag validation
- **utils/**: AWS-specific utility functions

#### OCM Integration (`pkg/ocm/`)
- **client/**: Generic OCM client interfaces using Go generics for type-safe resource operations
  - `SingleClusterSubResource[T]`: For 1:1 cluster resources (KubeletConfig, ClusterAutoscaler)
  - `CollectionClusterSubResource[T, S]`: For 1:many cluster resources (MachinePools, NodePools)
- **config/**: OCM configuration and token management
- **connection-builder/**: OCM connection establishment utilities

#### Validation Framework (`pkg/*/validations/`)
- Distributed validation utilities across multiple packages
- Common patterns for resource validation (cluster nodes, passwords, KMS ARNs)
- Reusable validation helpers and error handling

#### Utility Components
- **utils/parser/**: Parsing utilities including SQL parser, state machines, and string parsing
- **test/**: Testing utilities for VPC operations, KMS key management, and mock interfaces
- **log/**: Centralized logging utilities
- **deprecation/**: HTTP transport wrapper for handling deprecation headers

### Key Design Patterns

#### Generic Client Interfaces
The OCM client package uses Go generics to provide type-safe, reusable client interfaces:
```go
type SingleClusterSubResource[T any] interface {
    Get(ctx context.Context, clusterId string) (*T, error)
    Create(ctx context.Context, clusterId string, instance *T) (*T, error)
    Update(ctx context.Context, clusterId string, instance *T) (*T, error)
    Delete(ctx context.Context, clusterId string) error
}
```

#### AWS Client Factory
Centralized AWS client creation with support for multiple credential sources:
- Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
- Shared credentials file
- AWS profiles
- IAM roles

## Development Standards

### Commit Message Format
Follow conventional commits with JIRA ticket integration:
```
<type>[JIRA-TICKET] | [TYPE]: <MESSAGE>

[optional BODY]

[optional FOOTER(s)]
```

Types: `feat`, `fix`, `build`, `ci`, `docs`, `perf`, `refactor`, `style`, `test`

### Code Requirements
- All code must be covered by tests using Ginkgo
- Follow existing code patterns and conventions
- Use existing AWS SDK v2 and OCM SDK patterns
- Maintain backward compatibility for shared utilities

### JIRA Integration
- All changes require JIRA tickets in the OCM project
- Link PRs with `JIRA: OCM-xxxx` format
- Follow sprint workflow: Todo → In Progress → Code Review → Review → Done

## Testing Strategy

### Mock Generation
- Use `go.uber.org/mock` for mock generation
- Mock files located in `*/test/` directories
- OCM client mocks available in `pkg/ocm/client/test/`

### Test Organization
- Unit tests alongside source code (`*_test.go`)
- Test suites for package-level setup (`*_suite_test.go`)
- Integration test utilities in dedicated `test/` packages

## Common Patterns

### Error Handling
- Use structured logging via `pkg/log`
- Return detailed error information for AWS operations
- Handle HTTP status codes appropriately (404 for non-existent resources)

### Resource Management
- Context-aware operations throughout
- Proper cleanup of AWS resources in tests
- Pagination support for list operations

### Configuration
- Support multiple AWS credential methods
- Environment variable configuration where appropriate
- Default values defined in constants packages