# Development guide

## Table of contents

- [Requirements](#requirements)
- [Setup](#setup)
- [Run locally](#run-locally)
- [Build](#build)
- [Test](#test)
- [Vet and lint](#vet-and-lint)
- [CI](#ci)
- [Project structure](#project-structure)
- [Adding a new flag](#adding-a-new-flag)
- [Testing documentation examples](#testing-documentation-examples)

This guide covers local development commands for `treels`.

## Requirements

- Go `1.25.0` or compatible
- Optional: `golangci-lint` for local linting

## Setup

```bash
git clone https://github.com/OussamaM1/treels.git
cd treels
go mod download
```

## Run locally

```bash
go run .
go run . --tree --depth 2 --no-icons
go run . --long --readable
go run . --include "*.go" --exclude "*_test.go"
go run . --tree --git-status --no-icons
```

## Build

```bash
go build .
./treels --version
```

## Test

```bash
go test ./...
```

With coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Vet and lint

```bash
go vet ./...
golangci-lint run
```

The repository includes `.golangci.yml` with these enabled linters:

- `errcheck`
- `revive`
- `govet`
- `staticcheck`

## CI

GitHub Actions run:

- `go mod tidy`
- `go mod verify`
- `go vet ./...`
- `go build -v ./...`
- `go test -v ./... -coverprofile=coverage.out`
- `golangci-lint`

Workflow files live in `.github/workflows/`.

## Project structure

```text
cmd/       CLI command and flag definitions
module/    shared flag/options types, icons, colors, extension constants
service/   traversal, rendering, JSON output, gitignore matching
utils/     validation helpers
docs/      documentation
```

## Adding a new flag

When adding a user-facing flag, update:

1. `module/types.go`
2. `cmd/flag.go`
3. relevant service code
4. tests in `cmd/` and/or `service/`
5. `README.md` and relevant `docs/*.md`

## Testing documentation examples

Documentation examples should prefer `--no-icons` for stable text rendering unless the example is specifically about icons.
