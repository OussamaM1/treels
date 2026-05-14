---
name: golang-cli
description: >
  Guidance for building, extending, and reviewing this Go Cobra CLI. Use when
  changing command wiring, flags, filesystem traversal, terminal rendering,
  service-layer errors, or treels-style cmd/service/module/utils boundaries.
---

# Go CLI Development Rules

These rules are tailored to `github.com/oussamaM1/treels`, a Go 1.24 Cobra CLI
that combines `tree` and `ls` behavior with colored and icon-aware terminal
output.

## Architecture

Keep the package boundaries simple:

```text
main.go      calls cmd.Execute() only
cmd/         Cobra command construction, args, flags, process exit
service/     directory validation, traversal, counting, rendering orchestration
module/      plain data types and constants
utils/       pure helpers shared by cmd/service
```

Rules:
- `main.go` stays tiny and only calls `cmd.Execute()`.
- `cmd/` may import `service`, `module`, and Cobra.
- `service/` may import `module` and `utils`; it must not import `cmd`.
- `module/` must not depend on `cmd` or `service`.
- Do not add package-level mutable command state unless there is a strong reason.

## Cobra Commands

Use a command factory so tests get isolated command instances:

```go
func Execute() {
    if err := newRootCmd().Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func newRootCmd() *cobra.Command {
    var flags module.Flags
    cmd := &cobra.Command{
        Use:   "treels [path]",
        Args:  cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            options := module.Options{Flags: flags}
            if len(args) == 1 {
                options.Directory = args[0]
            }
            return service.Dispatcher(options)
        },
    }
    FlagDefinition(cmd, &flags)
    cmd.SilenceUsage = true
    return cmd
}
```

Rules:
- Use `RunE`, not `Run`, for commands that can fail.
- Register flags against the command instance passed to `FlagDefinition`.
- Keep filesystem validation in `service`, not `cmd`.
- Let `cmd.Execute()` be the only place that calls `os.Exit`.

## Service Errors

Service code must return errors instead of terminating the process.

Rules:
- Never call `log.Fatal`, `panic`, or `os.Exit` from `service/`, `utils/`, or
  `module/`.
- Validate and default the directory before printing any output.
- Wrap filesystem errors with the path involved, for example
  `fmt.Errorf("read directory %q: %w", path, err)`.
- Close opened directories on both success and failure paths.
- Return partial tree/list output only when that behavior is intentional and
  covered by tests.

## Rendering

Rendering should be testable without replacing global stdout.

Rules:
- Public CLI entry points may write to `os.Stdout`.
- Internal service rendering helpers should accept an `io.Writer`.
- Tests should pass a `bytes.Buffer` and assert on stable substrings.
- Do not assert terminal width, ANSI color escapes, or Nerd Font glyphs unless
  the test controls those conditions.

## Filesystem Traversal

Rules:
- Use `filepath.Join`, `filepath.Base`, and other path-aware stdlib helpers.
- Use `t.TempDir()` in tests; do not depend on the repository tree or `/tmp`.
- Hidden files are names whose first byte is `'.'`.
- Keep directory counting and file counting behavior explicit in tests whenever
  changing traversal.

## CI Contract

Before considering CLI changes ready, run the same checks as `.github/workflows`
expects:

```bash
go mod tidy
go mod verify
go vet ./...
go build ./...
go test ./...
```

Also run `go test -race ./...` locally for service or command changes. If
dependencies change, ensure `go mod tidy` leaves `go.mod` and `go.sum` in the
intended state.
