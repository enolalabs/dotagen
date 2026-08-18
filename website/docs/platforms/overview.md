---
sidebar_position: 1
---

# Platforms

All platform paths are relative to the current user's home directory. Agent
links point to files under `~/.dotagen/.generated/<target>/`; skill links point

| Platform | Target ID | Agent link | Skill link |
| --- | --- | --- | --- |
| Antigravity | `antigravity` | `~/.agents/<name>.md` | `~/.agent/skills/<name>/` |
| Claude Code | `claude-code` | `~/.claude/agents/<name>.md` | `~/.claude/skills/<name>/` |
| Codex | `codex` | `~/.codex/agents/<name>.md` | `~/.agents/skills/<name>/` |
| Cursor | `cursor` | `~/.cursor/rules/<name>.mdc` | `~/.cursor/skills/<name>/` |
| Gemini CLI | `gemini-cli` | `~/.gemini/agents/<name>.md` | `~/.gemini/skills/<name>/` |
| GitHub Copilot | `github-copilot` | `~/.github/agents/<name>.md` | `~/.github/skills/<name>/` |
| OpenCode | `opencode` | `~/.config/opencode/agents/<name>.md` | `~/.opencode/skills/<name>/` |
| Windsurf | `windsurf` | `~/.windsurf/rules/<name>.md` | `~/.windsurf/skills/<name>/` |

See the individual platform pages for the generated source paths and rendering
behavior.
