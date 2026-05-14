---
name: golang-unit-testing
description: >
  Guidance for deterministic Go tests in this repo. Use when adding or
  reviewing tests for Cobra commands, filesystem traversal, service rendering,
  utility validation, or CI behavior.
---

# Go Unit Testing Rules

These rules target the `treels` codebase: Go 1.24, Cobra, and packages split
across `cmd`, `service`, `module`, and `utils`.

## Test Placement

Rules:
- Put tests beside the code under test.
- Prefer the same package when testing unexported helpers such as `newRootCmd`,
  `dispatcher`, `isHidden`, or rendering helpers.
- Use black-box `package foo_test` only when intentionally testing the public API.

## Table Tests

Use named subtests for behavior matrices:

```go
func TestIsHidden(t *testing.T) {
    tests := []struct {
        name string
        input string
        want bool
    }{
        {name: "dot file", input: ".gitignore", want: true},
        {name: "regular file", input: "main.go", want: false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := isHidden(tt.input); got != tt.want {
                t.Fatalf("isHidden(%q) = %v, want %v", tt.input, got, tt.want)
            }
        })
    }
}
```

Rules:
- Case names describe behavior, not case numbers.
- Use `t.Fatalf` when later assertions would be meaningless.
- Use `t.Errorf` when a subtest can report several independent mismatches.

## Filesystem Tests

Rules:
- Always use `t.TempDir()` for fixture roots.
- Create files and directories in-process with helpers that call `t.Helper()`.
- Do not rely on permissions tests unless they are skipped or guarded for root
  and platform differences.
- Avoid tests that depend on the developer's current working directory.

## Command Tests

Rules:
- Test Cobra commands through `newRootCmd()` so flags and args are isolated per
  test.
- Use `cmd.SetArgs(...)`; do not mutate `os.Args`.
- Cover invalid paths, too many args, and at least one successful flag path when
  command behavior changes.
- Keep command tests focused on command wiring; service output details belong in
  `service` tests.

## Output Tests

Rules:
- Prefer `io.Writer` injection and `bytes.Buffer` over replacing `os.Stdout`.
- Assert stable substrings such as file names, connectors, and summary counts.
- Do not assert exact spacing from grid output unless terminal width is
  controlled by the test.
- Avoid assertions on ANSI colors or icons unless the test disables icons or
  controls the rendering mode.

## Error Tests

Every changed error path should have a test.

Cover:
- missing path
- path is a regular file
- too many CLI args
- read/open errors where deterministic
- no partial output before validation failures

## CI Contract

The GitHub workflows currently run:

```bash
go mod tidy
go mod verify
go vet ./...
go build -v ./...
go test -v ./...
golangci-lint
```

Before finishing test work, run at least:

```bash
go test ./...
go vet ./...
go test -race ./...
```

For code intended to merge, also run:

```bash
go mod tidy
go mod verify
go build ./...
```

If `golangci-lint` is installed locally, run `golangci-lint run`; otherwise say
that local lint verification was not available.
