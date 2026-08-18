---
sidebar_position: 4
---

# Codex

Target ID: `codex`.

- Agent output: `~/.dotagen/.generated/codex/<name>.md`
- Agent link: `~/.codex/agents/<name>.md`
- Skill output: `~/.dotagen/.generated/codex/skills/<name>/`
- Skill link: `~/.agents/skills/<name>/`

Both agents and skills render as Markdown with `name` and optional
`description` frontmatter. Note the intentional distinct skill link location:
the current implementation uses `~/.agents/skills/`, not `~/.codex/skills/`.
