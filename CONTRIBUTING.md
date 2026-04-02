# Contributing to s3tool

Thanks for your interest in contributing.

## Prerequisites

- Go 1.25+
- GNU Make
- Optional: OpenTofu (`tofu`) for deployment test helpers
- Optional: Docker (for local MinIO testing)

## Setup

```bash
git clone https://github.com/schidstorm/s3tool.git
cd s3tool
go mod download
```

## Run locally

```bash
go run ./cmd/s3tool
```

Or with the provided make target:

```bash
make debug
```

## Run tests

```bash
make test
make test-cover
make lint
make vet
```

This runs unit tests with coverage and generates `coverage.html`.

## Useful development commands

```bash
make build
make fmt
make generate-screens
make create_test_bucket
make delete_test_bucket
```

For local object storage integration testing (MinIO), see `docs/INTEGRATION_TESTING.md`.

## Coding guidelines

- Keep changes small and focused.
- Add tests for new behavior and bug fixes.
- Preserve current package structure and naming conventions.
- Ensure `go test ./...` passes before opening a PR.

## Commit messages

This project uses Conventional Commits.

Examples:

- `feat: add profile sorting in profile page`
- `fix: handle empty bucket list`
- `test: add object page navigation tests`

## Pull requests

- Explain the problem and the proposed solution.
- Link related issue(s).
- Mention any breaking changes explicitly.

## Reporting bugs

Open a GitHub issue with:

- Environment details (OS, Go version)
- Reproduction steps
- Expected vs actual behavior
