---
sidebar_position: 3
---

# Rendering And Symlinks

dotagen does not write source definitions directly into platform folders. It
uses this flow:

```text
~/.dotagen/agents/*.md and ~/.dotagen/skills/*/SKILL.md
  -> platform renderer
  -> ~/.dotagen/.generated/<platform>/...
  -> symlink in the platform's home-directory path
```

Agent output is a generated file. Skill output is a generated directory that
contains `SKILL.md` and copied reference files. The platform-facing file or
directory is a symlink to the generated output.

For an agent named `da-review` sent to Claude Code, the generated file is
`~/.dotagen/.generated/claude-code/da-review.md` and the linked platform file
is `~/.claude/agents/da-review.md`.

Sync does not overwrite a non-symlink platform file. If a destination exists
and is not a symlink, sync fails rather than replacing it. `dotagen status`
reports synced, out-of-date, missing, and broken symlink states; run
`dotagen sync` after source or configuration changes.

`dotagen clean` removes dotagen-managed agent and skill symlinks for the
current user and empties `~/.dotagen/.generated/`. It does not delete source
definitions or change `config.yaml`.
