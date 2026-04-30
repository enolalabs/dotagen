# dotagen

> **Define sub-agents once, inject everywhere.**

`dotagen` is a Go CLI tool that lets you define coding sub-agents **once** in Markdown and automatically distribute them across multiple coding agent platforms (Claude Code, Cursor, Gemini CLI, OpenCode).

Instead of writing and maintaining N sets of configurations for N platforms, you manage **a single source of truth** in `.dotagen/` — dotagen renders each agent to the correct platform format and creates symlinks automatically.

## Why dotagen?

| Problem | Solution |
|---|---|
| Rewrite the same instructions for every agent platform | Define once in Markdown, dotagen renders to each format |
| Maintain N config sets for N platforms | Centralized management in `.dotagen/agents/` |
| Re-setup every time you switch tools | Run `dotagen sync` — all platforms updated |
| No visibility into which agents are active where | `dotagen status` shows detailed sync state |

## Supported Platforms

| Platform | Output Path | Format |
|---|---|---|
| **Claude Code** | `.claude/agents/{name}.md` | Pure Markdown |
| **Cursor** | `.cursor/rules/{name}.mdc` | YAML frontmatter + Markdown |
| **Gemini CLI** | `.gemini/agents/{name}.md` | Pure Markdown |
| **OpenCode** | `.opencode/agents/{name}.md` | YAML frontmatter + Markdown |

## Installation

### Build from source

```bash
git clone https://github.com/k0walski/dotagen.git
cd dotagen
make build
```

The binary will be created at `./dotagen`.

### Or install directly

```bash
go install github.com/k0walski/dotagen/cmd/dotagen@latest
```

## Built-in Agents

dotagen ships with **144 built-in agents** covering a wide range of specialties. They are injected automatically when you run `dotagen init`, so you don't need to create agents from scratch.

| Category | Examples |
|---|---|
| Core Development | `backend-developer`, `frontend-developer`, `fullstack-developer`, `api-designer`, `ui-designer` |
| Language Specialists | `golang-pro`, `rust-engineer`, `python-pro`, `typescript-pro`, `nextjs-developer` |
| Infrastructure | `devops-engineer`, `docker-expert`, `kubernetes-specialist`, `terraform-engineer` |
| Quality & Security | `code-reviewer`, `qa-expert`, `penetration-tester`, `security-auditor` |
| Data & AI | `ai-engineer`, `data-engineer`, `llm-architect`, `ml-engineer`, `prompt-engineer` |
| Developer Experience | `build-engineer`, `cli-developer`, `documentation-engineer`, `refactoring-specialist` |
| Specialized Domains | `blockchain-developer`, `game-developer`, `fintech-engineer`, `iot-engineer` |
| Business & Product | `product-manager`, `business-analyst`, `technical-writer`, `ux-researcher` |
| Meta-Orchestration | `agent-organizer`, `multi-agent-coordinator`, `task-distributor`, `workflow-orchestrator` |
| Research & Analysis | `market-researcher`, `competitive-analyst`, `data-researcher`, `trend-analyst` |

All agents are **disabled by default**. You decide which agents to enable and for which platforms.

## Quick Start

### 1. Initialize

```bash
dotagen init
```

Creates `.dotagen/` with all 144 built-in agents and a config file where every agent has empty targets (disabled):

```
.dotagen/
├── config.yaml       # Configuration — set targets to enable agents
├── agents/           # 144 built-in agent definitions (*.md)
├── .generated/       # Rendered output (git-ignored)
└── .gitignore
```

### 2. Configure

Edit `.dotagen/config.yaml` to enable agents by setting their `targets`. By default all agents have `targets: []` (disabled):

```yaml
targets:
  - claude-code
  - cursor
  - gemini-cli
  - opencode

agents:
  backend-developer:
    targets: all              # Enable for all platforms
  code-reviewer:
    targets: all
  frontend-developer:
    targets:
      - cursor                # Cursor only
      - opencode              # OpenCode only
  golang-pro:
    targets:
      - claude-code           # Claude Code only
  # ... all other agents remain disabled (targets: [])
```

**Agent `targets` values:**
- `all` — Apply to all configured platforms
- A specific list, e.g. `[claude-code, cursor]`
- `[]` — Disabled (not synced)

### 3. Sync

```bash
# Sync all platforms
dotagen sync

# Sync a single platform
dotagen sync cursor
```

Output:

```
✓ Synced 3 agents to 4 platform(s)

  claude-code:
    ✓ .claude/agents/review-code.md → .dotagen/.generated/claude-code/review-code.md
    ✓ .claude/agents/testing.md → .dotagen/.generated/claude-code/testing.md

  cursor:
    ✓ .cursor/rules/review-code.mdc → .dotagen/.generated/cursor/review-code.mdc
    ✓ .cursor/rules/planning.mdc → .dotagen/.generated/cursor/planning.mdc
    ✓ .cursor/rules/testing.mdc → .dotagen/.generated/cursor/testing.mdc

  ...
```

### 4. Check Status

```bash
dotagen status
```

Shows the state of each agent on each platform:

- `✓ synced` — Up to date
- `⚠ out-of-date` — Source has changed, needs re-sync
- `✗ missing` — Not yet created
- `💔 broken` — Symlink is broken

## CLI Commands

| Command | Description |
|---|---|
| `dotagen init` | Initialize `.dotagen/` with 144 built-in agents (all disabled) |
| `dotagen sync [target]` | Render & symlink agents. Specify a `target` to sync only that platform (e.g. `cursor`) |
| `dotagen status` | Show sync status of all agents |
| `dotagen clean` | Remove all generated files and symlinks |
| `dotagen serve` | Start web dashboard at `http://localhost:7890` |
| `dotagen --version` | Print version |
| `dotagen --help` | Print help |

## Web Dashboard

```bash
dotagen serve
```

Starts a web dashboard at `http://localhost:7890` with the following features:

- **Agent Management** — CRUD agents directly from the UI
- **Target Matrix** — Assign/unassign agents to platforms via toggle
- **Preview** — View rendered output for each platform
- **Sync/Clean** — Trigger sync or clean from the web UI
- **Status** — View status of all symlinks

### REST API

The dashboard is powered by a REST API you can also call directly:

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/config` | Get current config |
| `PUT` | `/api/config` | Update config |
| `GET` | `/api/agents` | List all agents |
| `GET` | `/api/agents/{name}` | Get agent detail |
| `POST` | `/api/agents` | Create a new agent |
| `PUT` | `/api/agents/{name}` | Update an agent |
| `DELETE` | `/api/agents/{name}` | Delete an agent |
| `GET` | `/api/targets` | List available targets |
| `GET` | `/api/preview/{agent}/{target}` | Preview rendered output |
| `POST` | `/api/sync` | Trigger full sync |
| `POST` | `/api/sync/{target}` | Trigger sync for one target |
| `POST` | `/api/clean` | Trigger clean |
| `GET` | `/api/status` | Get overall status |

### Options

```bash
dotagen serve --port 8080     # Custom port (default: 7890)
dotagen serve --open=false    # Don't auto-open browser
```

## Config Reference

### `.dotagen/config.yaml`

```yaml
# Platforms to target
targets:
  - claude-code
  - cursor
  - gemini-cli
  - opencode

# Agent-to-target mapping
agents:
  <agent-name>:
    targets: all | [target1, target2, ...]
    disabled: true | false   # (optional) Temporarily disable an agent
```

### Agent Markdown Format

```markdown
---
description: Short description of the agent
category: classification      # optional
---

# Agent Name

Agent instructions written in Markdown.

## Guidelines

- Rule 1
- Rule 2
```

**Note:** Frontmatter (`---`) is optional. When present, the `description` field is used by the Cursor and OpenCode adapters when rendering output.

## Architecture

```
.dotagen/agents/*.md  +  .dotagen/config.yaml
         │
         ▼
    Config Parser → Agent Parser → Template Renderer
                                        │
                                        ▼
                              .dotagen/.generated/{platform}/
                                        │
                                        ▼
                              Symlink Manager → .claude/agents/, .cursor/rules/, ...
```

## Development

### Requirements

- Go 1.26+

### Commands

```bash
make build     # Build binary
make test      # Run tests
make lint      # Run linter
make dev       # Run development server
make install   # Install to $GOPATH/bin
make clean     # Remove build artifacts
```

### Project Structure

```
├── cmd/dotagen/main.go          # Entry point
├── internal/
│   ├── agent/                   # Agent markdown parser
│   ├── builtin/                 # Built-in agents (go:embed)
│   │   ├── embed.go
│   │   └── agents/              # 144 agent .md files
│   ├── cli/                     # CLI commands (cobra)
│   ├── config/                  # Config parser & validator
│   ├── engine/                  # Renderer & symlink manager
│   ├── platform/                # Platform adapters
│   │   ├── adapter.go           # Adapter interface
│   │   ├── claude_code.go
│   │   ├── cursor.go
│   │   ├── gemini_cli.go
│   │   ├── opencode.go
│   │   └── registry.go
│   └── web/                     # Web dashboard (go:embed)
│       ├── server.go
│       ├── api.go
│       └── static/
├── go.mod
├── Makefile
└── README.md
```

## License

MIT
