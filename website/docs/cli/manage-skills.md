---
sidebar_position: 4
---

# Manage Skills

## List installed skills

```bash
dotagen skill list
```

Lists parsed skills from `~/.dotagen/skills/` with category, resolved targets,
and reference-file count.

## Create a skill

```bash
dotagen skill create release-check
dotagen skill create release-check --targets claude-code,opencode
```

The command adds the `ds-` prefix when omitted, creates
`~/.dotagen/skills/ds-release-check/SKILL.md`, and registers the skill. Names
must use lowercase letters, numbers, hyphens, or underscores. `--targets`
defaults to `all`.

## Delete a skill

```bash
dotagen skill delete release-check
dotagen skill delete release-check --force
```

The command adds the `ds-` prefix when omitted. Without `--force`, it prints a
destructive-operation notice and makes no change. With `--force`, it removes
the skill directory and its configuration entry. Sync or clean afterward to
remove any existing platform links.
