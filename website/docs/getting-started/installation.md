---
sidebar_position: 1
---

# Installation

## Release installer

On Linux or macOS, install the current release with:

```bash
curl -fsSL https://raw.githubusercontent.com/enolalabs/dotagen/main/install.sh | sh
```

The installer chooses the matching release binary and installs it to
`/usr/local/bin`, or `~/.local/bin` when that is the available fallback.

## Go install

With Go 1.26 or newer:

```bash
go install github.com/enolalabs/dotagen/cmd/dotagen@latest
```

## Build from source

```bash
git clone https://github.com/enolalabs/dotagen.git
cd dotagen
make build
```

Put the resulting `dotagen` binary on your `PATH`.

## Verify the installation

Use the `version` subcommand:

```bash
dotagen version
dotagen help
```

`dotagen --version` is not a documented command. The CLI defines `dotagen
version`; older README examples using `--version` are stale.

Continue with [initialization](./initialize.md).
