# Contributing

## Table of contents

- [Getting started](#getting-started)
- [Development workflow](#development-workflow)
- [Contribution guidelines](#contribution-guidelines)
- [Adding or changing CLI behavior](#adding-or-changing-cli-behavior)
- [Documentation style](#documentation-style)
- [Reporting bugs](#reporting-bugs)
- [Feature requests](#feature-requests)

Thanks for your interest in contributing to `treels`.

## Getting started

```bash
git clone https://github.com/OussamaM1/treels.git
cd treels
go test ./...
```

## Development workflow

Before opening a pull request, run:

```bash
go test ./...
go vet ./...
golangci-lint run
```

If you do not have `golangci-lint` installed, CI will still run it on the pull request.

## Contribution guidelines

- Keep changes focused and easy to review.
- Add or update tests for behavior changes.
- Update documentation for user-facing changes.
- Prefer existing project patterns before adding new abstractions.
- Avoid introducing new dependencies unless they are clearly justified.

## Adding or changing CLI behavior

For a new flag or output behavior, update:

- `cmd/flag.go`
- `module/types.go`
- relevant implementation under `service/`
- tests under `cmd/` and/or `service/`
- `README.md`
- relevant files under `docs/`

## Documentation style

- Prefer examples using `--no-icons` unless documenting icon behavior.
- Keep README concise and link to detailed docs.
- Document flag interactions when behavior may be surprising.
- Be explicit about limitations, especially for `.gitignore` and JSON output.

## Reporting bugs

When reporting a bug, include:

- operating system
- `treels --version`
- command used
- expected output
- actual output
- small reproduction directory structure, if possible

## Feature requests

Good feature requests explain:

- the problem or workflow
- proposed CLI usage
- expected output
- any interaction with existing flags
