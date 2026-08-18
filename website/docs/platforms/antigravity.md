---
sidebar_position: 2
---

# Antigravity

Target ID: `antigravity`.

- Agent output: `~/.dotagen/.generated/antigravity/<name>.md`
- Agent link: `~/.agents/<name>.md`
- Skill output: `~/.dotagen/.generated/antigravity/skills/<name>/`
- Skill link: `~/.agent/skills/<name>/`

Agents render as plain Markdown. Skills render with `name` and optional
`description` frontmatter.

When Antigravity is selected during sync, dotagen also writes every parsed
agent as a global workflow at
`~/.gemini/antigravity/global_workflows/<name>.md`. These workflow files use
agent content plus description frontmatter, are not symlinks, and stale files
in that directory are removed by the workflow sync.
