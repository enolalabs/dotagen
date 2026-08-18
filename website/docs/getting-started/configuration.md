---
sidebar_position: 3
---

# Configuration

The only configuration file is `~/.dotagen/config.yaml`. It selects the
platforms available to the store and the platforms enabled for each agent or
skill.

```yaml
targets:
  - claude-code
  - gemini-cli
  - opencode

agents:
  da-review:
    targets: [claude-code, opencode]
  da-disabled:
    targets: all
    disabled: true

skills:
  ds-release-check:
    targets: all
  ds-gemini-only:
    targets: gemini-cli
```

The top-level `targets` field is a YAML sequence. At an agent or skill entry,
`targets` accepts either a YAML string or a YAML sequence. `all` resolves to
the top-level `targets` list. An empty list disables that entry. `disabled:
true` also prevents the entry from resolving to a target.

The supported selectors are `antigravity`, `claude-code`, `codex`, `cursor`,
`gemini-cli`, `github-copilot`, `opencode`, and `windsurf`. Use only these
values and `all` for per-entry target selection.

## Sync selection precedence

`dotagen sync` without an argument first detects existing platform agent
directories under your home directory. When it finds one or more, it syncs to
those detected targets instead of the configured top-level `targets` list. Only
when it detects none does it fall back to `config.yaml`.

To bypass that selection and sync one known target, use `dotagen sync
<target>`, for example `dotagen sync opencode`.
