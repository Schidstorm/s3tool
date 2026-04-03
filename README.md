# s3tool

![CI](https://github.com/schidstorm/s3tool/actions/workflows/go.yml/badge.svg)
![Go Version](https://img.shields.io/github/go-mod/go-version/schidstorm/s3tool)
![Release](https://img.shields.io/github/v/release/schidstorm/s3tool)

Terminal-based S3 client for browsing buckets and objects across AWS profiles and S3-compatible endpoints.

## Features

- Interactive terminal UI for profile, bucket, and object navigation
- AWS profile discovery from `~/.aws/config` and `~/.aws/credentials`
- Custom S3 profile loading from YAML files in `~/.s3tool` (configurable)
- Support for S3-compatible endpoints (for example MinIO)
- Shell completion generation via Cobra (`bash`, `zsh`, `fish`, `powershell`)

## Installation

### Option 1: Install with Go

```bash
go install github.com/schidstorm/s3tool/cmd/s3tool@latest
```

### Option 2: Build locally

```bash
git clone https://github.com/schidstorm/s3tool.git
cd s3tool
go build -o build/s3tool ./cmd/s3tool
./build/s3tool
```

## Quick Start

### Use existing AWS profiles

`s3tool` auto-loads profiles from your AWS CLI config/credentials files.

```bash
s3tool
```

### Add a custom S3-compatible profile

Create a file like `~/.s3tool/minio.yaml`:

```yaml
access_key_id: minioadmin
secret_access_key: minioadmin
base_endpoint: http://localhost:9000
use_path_style: true
region: us-east-1
```

Then run:

```bash
s3tool
```

### CLI options

```bash
s3tool --help
s3tool completion --shell zsh
```

Main flags:

- `-p, --profiles`: path to directory containing profile YAML files
- `--loaders.aws`: enable AWS profile loader (default: true)
- `--loaders.s3tool`: enable YAML profile loader (default: true)
- `--loaders.memory`: test-only in-memory loader (hidden)

## Development

Common commands:

```bash
make test
make test-cover
make lint
make vet
make fmt
make build
make debug
make generate-screens
```

Integration test setup with MinIO is documented in `docs/INTEGRATION_TESTING.md`.

## Compatibility

See `docs/COMPATIBILITY.md` for supported Go versions and S3-compatible providers.

## Documentation

- Architecture: `docs/ARCHITECTURE.md`
- Command reference: `docs/COMMANDS.md`
- Roadmap: `docs/ROADMAP.md`
- Governance: `docs/GOVERNANCE.md`
- Release process: `docs/RELEASE_PROCESS.md`
- Telemetry policy: `docs/TELEMETRY.md`

## Roadmap

See `TODOS.md` and `docs/ROADMAP.md` for open-source hardening and roadmap items.

## License

This project is licensed under the MIT License. See `LICENSE` for details.
