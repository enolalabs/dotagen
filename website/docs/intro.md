---
sidebar_position: 2
---

# Introduction

dotagen is a Go CLI for maintaining coding sub-agents and skills once and
deploying them to multiple local coding-agent platforms. Its source of truth is
global to your user account, not a repository:

```text
~/.dotagen/
├── agents/           # One Markdown file per agent
├── skills/           # One directory per skill, with SKILL.md
├── .generated/       # Platform-specific rendered output
└── config.yaml       # Target selection for agents and skills
```

Run `dotagen init` to create the store and import the bundled material. Enable
the definitions you want in `~/.dotagen/config.yaml`, then run `dotagen sync`.
dotagen writes rendered files below `.generated/` and creates symlinks in the
platform directories under your home directory.

Supported target IDs are `antigravity`, `claude-code`, `codex`, `cursor`,
`gemini-cli`, `github-copilot`, `opencode`, and `windsurf`.

## What to read next

- [Install dotagen](./getting-started/installation.md)
- [Initialize the global store](./getting-started/initialize.md)
- [Configure targets](./getting-started/configuration.md)
- [Understand rendering and symlinks](./concepts/rendering-and-symlinks.md)
