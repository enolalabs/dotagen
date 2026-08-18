---
sidebar_position: 2
---

# Init, Sync, And Status

## `dotagen init`

Creates `~/.dotagen/` with `agents/`, `skills/`, `.generated/`,
`.gitignore`, and `config.yaml`. It imports built-in content. When the store
already exists, it asks before overwriting it; see [Initialize](../getting-started/initialize.md)
for what overwrite removes.

## `dotagen sync [target]`

Renders configured agents and skills, creates platform directories as needed,
and links the generated output into them.

```bash
dotagen sync
dotagen sync claude-code
dotagen sync windsurf
```

The optional target must be one of the eight target IDs. Without an argument,
existing platform folders under the user's home directory take precedence over
the configured top-level targets. Sync also removes stale dotagen-managed links
for the selected targets and writes Antigravity global workflows when
`antigravity` is selected.

## `dotagen status`

Shows every parsed agent and skill against each configured top-level target.
Agent states are `synced`, `out of date`, `missing`, or `broken symlink`.
Skill state reports whether its linked directory and `SKILL.md` are present.

## `dotagen clean`

Removes managed `da-` agent links and managed skill links, then empties
`~/.dotagen/.generated/`. It is destructive to generated output and platform
links, but it does not ask for confirmation and leaves sources and configuration
in place. Run `dotagen sync` to recreate the links.
