---
sidebar_position: 3
---

# Create Agents

Create a named agent and register it in `config.yaml`:

```bash
dotagen create review-code
dotagen create review-code --description "Senior code reviewer" --targets claude-code,gemini-cli
dotagen create release-check --content "# Release Check"
dotagen create release-check --file ./agent-template.md
dotagen create release-check --template
```

Agent names must start with an alphanumeric character and otherwise use only
letters, digits, hyphens, or underscores. The command writes
`~/.dotagen/agents/<name>.md`.

Content priority is `--content`, then `--file`, then `--template`. With none
of those flags, dotagen opens `$EDITOR`, then `$VISUAL`, and otherwise reads
the agent Markdown from standard input. `--description` prepends a
`description` frontmatter field.

`--targets` is a comma-separated list and defaults to `all`. Run `dotagen
sync` after creating the agent.
