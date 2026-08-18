---
sidebar_position: 5
---

# Serve, Update, And Uninstall

## `dotagen serve`

```bash
dotagen serve
dotagen serve --port 8080
```

Starts the embedded dashboard. The default port is `7890`. The CLI prints a
localhost URL, but the server actually binds to `:<port>`, which exposes it on
all network interfaces. The dashboard and REST API have no authentication.
Use it only on a trusted local network; do not expose the port to an untrusted
network until the code provides stronger binding or authentication.

## `dotagen update`

Checks GitHub Releases for `enolalabs/dotagen`, downloads the matching current
OS and architecture binary, replaces the running binary, then invokes its
`version` subcommand. It requires network access and a release asset matching
the platform naming convention.

## `dotagen uninstall [--force]`

Uninstall removes managed agent and skill links, Antigravity global workflow
files, empty dotagen-created platform directories, `~/.dotagen/`, and attempts

Without `--force`, it asks for the exact confirmation `yes`. Use `--force` only
when you intend to remove all of those artifacts without the prompt.
