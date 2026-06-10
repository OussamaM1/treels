# Gitignore support

Use `--gitignore` to hide entries matched by `.gitignore` rules from the target directory.

```bash
treels --tree --gitignore
treels --json --gitignore
treels --tree --gitignore --dirs-only
```

## Current behavior

`treels` reads `.gitignore` from the directory being listed.

```bash
treels --gitignore path/to/project
```

In this example, rules are read from:

```text
path/to/project/.gitignore
```

If that file does not exist, `--gitignore` behaves like normal listing.

## Supported rule features

The matcher supports common `.gitignore` pattern features:

| Feature | Example |
| --- | --- |
| Comments | `# generated files` |
| Escaped comments | `\#literal-file` |
| Negation | `!keep.log` |
| Anchored rules | `/dist` |
| Directory-only rules | `build/` |
| File globs | `*.log` |
| Slash-separated globs | `logs/*.txt` |
| `**` wildcard | `logs/**/*.tmp` |

## Examples

Given this `.gitignore`:

```gitignore
*.log
dist/
!important.log
```

This command:

```bash
treels --gitignore --no-icons
```

Will:

- hide files ending in `.log`
- hide the `dist` directory
- keep `important.log`

## Interaction with other flags

| Combination | Behavior |
| --- | --- |
| `--gitignore --tree` | Ignored directories are not recursively traversed. |
| `--gitignore --dirs-only` | Ignored directories are omitted; files are already omitted by `--dirs-only`. |
| `--gitignore --all` | Hidden files may be shown, but ignored hidden files are still omitted. |
| `--gitignore --json` | JSON entries and summary counts reflect filtered output. |

## Limitations

`treels` does not currently aim for full Git parity. Important limitations:

- Nested `.gitignore` files are not loaded.
- Global Git excludes are not loaded.
- `.git/info/exclude` is not loaded.
- Rules are interpreted relative to the target directory only.

These limitations are good candidates for future improvement.
