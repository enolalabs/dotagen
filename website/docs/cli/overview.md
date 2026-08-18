---
sidebar_position: 1
---

# CLI Overview

Run `dotagen help` to see Cobra-generated help, or `dotagen help <command>`
for a command's flags and arguments.

| Command | Purpose |
| --- | --- |
| `init` | Create or reinitialize `~/.dotagen` |
| `sync [target]` | Render and link enabled agents and skills |
| `status` | Report sync state |
| `clean` | Remove managed links and generated output |
| `create <name>` | Create and register an agent |
| `skill list/create/delete` | Manage skills |
| `serve --port <port>` | Start the dashboard |
| `update` | Replace the binary with the latest GitHub release |
| `uninstall --force` | Remove dotagen artifacts and binary |
| `version` | Print build and platform information |
| `completion <shell>` | Generate shell completion |

## Known command-documentation mismatches

- Use `dotagen version` to print version information.
- Use `dotagen serve --port <port>` to choose the dashboard port.
- The `sync` long help text names only five targets, but the command validates
  and accepts all eight current target IDs listed in the configuration
  reference.

See [init, sync, and status](./init-sync-status.md),
[agent creation](./create-agents.md), [skill management](./manage-skills.md),
and [serve, update, and uninstall](./serve-update-uninstall.md).
