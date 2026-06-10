# Icons and fonts

## Table of contents

- [Requirements](#requirements)
- [Disable icons](#disable-icons)
- [Colors](#colors)
- [File type support](#file-type-support)

`treels` uses Nerd Font icons by default to make file types easier to scan.

## Requirements

Icons require a terminal font that includes Nerd Font glyphs.

Recommended fonts:

- FiraCode Nerd Font
- JetBrainsMono Nerd Font
- Hack Nerd Font
- MesloLGS Nerd Font

After installing a Nerd Font, select it in your terminal profile.

## Disable icons

If icons render as empty boxes, question marks, or strange symbols, disable them:

```bash
treels --no-icons
```

This is also useful for:

- CI logs
- plain text snapshots
- terminals without glyph support
- documentation examples

## Colors

`treels` applies ANSI colors to icons and directory names in text output.

JSON output never includes icons or ANSI color codes.

## File type support

Icons are selected using filename and extension mappings. Examples include:

| Type | Examples |
| --- | --- |
| Go | `.go`, `.mod`, `.sum` |
| Web | `.html`, `.css`, `.js`, `.ts`, `.jsx`, `.tsx`, `.vue` |
| Data | `.json`, `.yml`, `.yaml`, `.xml`, `.sql` |
| Media | `.png`, `.jpg`, `.mp4`, `.mp3` |
| Archives | `.zip`, `.tar`, `.gz`, `.rar`, `.7z` |
| Project files | `README.md`, `LICENSE`, `Dockerfile`, `package.json`, `Cargo.toml` |

Unknown files use the default file icon. Unknown directories use the default folder icon.
