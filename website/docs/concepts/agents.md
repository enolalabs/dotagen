---
sidebar_position: 1
---

# Agents

An agent is one Markdown file in `~/.dotagen/agents/`. Its filename, without
the `.md` suffix, is its identifier. For example,
`~/.dotagen/agents/da-review.md` defines `da-review`.

Markdown body content is the portable instruction set. Optional YAML
frontmatter adds metadata consumed by particular renderers:

```markdown
---
description: Review a change for correctness and regressions.
model: gpt-5
tools: read,search
globs: "src/**/*.go"
---

# Review Changes

Read the change, identify behavioral risks, and report findings first.
```

The parser tolerates missing or unparsable frontmatter and treats the file body
as the agent content. Individual platform renderers decide which metadata to
preserve. For example, Claude Code preserves `description` and `model`, Gemini
CLI also renders `tools`, Cursor uses `description` and `globs`, and OpenCode
sets `mode: subagent`.

Create an agent through the CLI or add the file directly, then add an entry
under `agents` in `config.yaml` and run `dotagen sync`.
