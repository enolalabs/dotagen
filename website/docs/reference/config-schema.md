---
sidebar_position: 1
---

# Configuration Schema

`~/.dotagen/config.yaml` has three top-level fields:

```yaml
targets:                 # Required target universe for all
  - antigravity
  - claude-code

agents:
  da-review:
    targets: all         # "all", one string, a list, or []
    disabled: false      # Optional; true disables the entry

skills:
  ds-review-checklist:
    targets:
      - claude-code
    disabled: false
```

| Field | Type | Meaning |
| --- | --- | --- |
| `targets` | string list | Configured platform IDs |
| `agents` | map | Agent ID to target selection and disabled state |
| `skills` | map | Skill ID to target selection and disabled state |
| `<entry>.targets` | string or string list | `all`, a selector, selectors, or empty list |
| `<entry>.disabled` | boolean | Optional explicit disable switch |

The valid platform IDs are `antigravity`, `claude-code`, `codex`, `cursor`,
`gemini-cli`, `github-copilot`, `opencode`, and `windsurf`. `all` at an entry
resolves to the top-level `targets` list. Unknown IDs in `targets` and agent
target lists fail CLI configuration validation; keep skill target lists to the
same valid IDs as well.

Each configured agent should have a source file at
`~/.dotagen/agents/<name>.md`. Each configured skill should have
`~/.dotagen/skills/<name>/SKILL.md`.
