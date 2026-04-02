# Open Source Readiness TODO

This checklist is prioritized to make `s3tool` easier to discover, use, and contribute to.

## P0 - Must Have

- [x] Add `README.md`
  - Project overview and motivation
  - Feature list and current status
  - Install instructions (`go install`, release binary)
  - Quick start with examples
  - Screenshot section using files from `screens/`
- [x] Add `LICENSE` (e.g., MIT or Apache-2.0)
- [x] Add `CONTRIBUTING.md`
  - Local setup steps
  - How to run tests (`make tests`)
  - Branching and PR expectations
- [x] Add `CODE_OF_CONDUCT.md`
- [x] Add `SECURITY.md`
  - Supported versions
  - Vulnerability reporting contact/process
- [ ] Add issue and PR templates in `.github/`
  - Bug report template
  - Feature request template
  - Pull request template

## P1 - CI/CD and Quality Gates

- [ ] Tighten CI in `.github/workflows/go.yml`
  - Add `go test -cover ./...` and enforce minimum coverage target
  - Add `go vet ./...`
  - Add `golangci-lint` run
  - Cache Go modules/build cache to speed up CI
- [x] Separate Renovate execution from build/test workflow
  - Keep CI deterministic for forks and external contributors
- [ ] Add GitHub release workflow
  - Build cross-platform binaries (Linux/macOS/Windows)
  - Attach checksums and release notes
- [ ] Add status badges to `README.md`
  - Build status
  - Go version
  - Latest release

## P2 - Developer Experience

- [ ] Add `make` targets for common contributor tasks
  - `lint`, `vet`, `fmt`, `test`, `test-cover`, `build`
- [ ] Add pinned tool versions (lint/tooling) for reproducibility
- [ ] Add `.editorconfig` and normalize formatting expectations
- [ ] Expand test coverage for terminal UI flows and connector edge cases
- [ ] Add integration test instructions for local object-storage emulation (MinIO)

## P3 - Documentation and Governance

- [ ] Add architecture documentation
  - Connector/loader model
  - Terminal page navigation and state flow
- [ ] Add command reference docs (generated from Cobra help)
- [ ] Add compatibility matrix
  - Supported Go version(s)
  - Supported S3-compatible providers
- [ ] Add roadmap and known limitations
- [ ] Add maintainership and review policy
  - Who can review/merge
  - Expected response times

## P4 - Trust and Community Signals

- [ ] Add project metadata
  - Topics and clear repo description
  - Website/docs link
- [x] Add Dependabot security updates or keep Renovate strictly configured for security
- [ ] Add signed tags/releases and changelog discipline for each release
- [ ] Add basic telemetry policy statement (if telemetry is introduced later)

## Suggested First Milestone (1-2 weeks)

- [x] `README.md`, `LICENSE`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md`
- [x] CI with build + test + lint + vet
- [ ] Issue/PR templates
- [ ] First tagged release with release notes
