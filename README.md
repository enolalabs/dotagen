# dotagen

> **Define sub-agents and skills once, inject everywhere.**

![Overview](docs/screenshots/overview.png)

`dotagen` is a Go CLI tool that lets you define coding sub-agents and skills **once** in Markdown and automatically distribute them across multiple coding agent platforms (Antigravity, Claude Code, Codex, Cursor, Gemini CLI, GitHub Copilot, OpenCode, Windsurf).

Instead of writing and maintaining N sets of configurations for N platforms, you manage **a single source of truth** in `.dotagen/` — dotagen renders each agent/skill to the correct platform format and creates symlinks automatically.

## Why dotagen?

| Problem | Solution |
|---|---|
| Rewrite the same instructions for every agent platform | Define once in Markdown, dotagen renders to each format |
| Maintain N config sets for N platforms | Centralized management in `.dotagen/` |
| Re-setup every time you switch tools | Run `dotagen sync` — all platforms updated |
| No visibility into which agents are active where | `dotagen status` shows detailed sync state |

## Supported Platforms

| Platform | ID | Agent Path | Skill Path |
|---|---|---|---|
| **Antigravity** | `antigravity` | `.agents/{name}.md` | `.agent/skills/{name}/` |
| **Claude Code** | `claude-code` | `.claude/agents/{name}.md` | `.claude/skills/{name}/` |
| **Codex** | `codex` | `.codex/agents/{name}.md` | `.agents/skills/{name}/` |
| **Cursor** | `cursor` | `.cursor/rules/{name}.mdc` | `.cursor/skills/{name}/` |
| **Gemini CLI** | `gemini-cli` | `.gemini/agents/{name}.md` | `.gemini/skills/{name}/` |
| **GitHub Copilot** | `github-copilot` | `.github/agents/{name}.md` | `.github/skills/{name}/` |
| **OpenCode** | `opencode` | `.config/opencode/agents/{name}.md` | `.opencode/skills/{name}/` |
| **Windsurf** | `windsurf` | `.windsurf/rules/{name}.md` | `.windsurf/skills/{name}/` |

## Installation

### One-liner (Linux & macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/enolalabs/dotagen/main/install.sh | sh
```

Or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/enolalabs/dotagen/main/install.sh | sh
```

Supports **Linux** (amd64, arm64) and **macOS** (Intel, Apple Silicon). Automatically detects your OS/arch, downloads the latest release, and installs to `/usr/local/bin` (or `~/.local/bin` as fallback).

### Update

Already installed? Run the same install command again, or use:

```bash
dotagen update
```

Checks for the latest version and self-updates in place.

### Windows (PowerShell)

```powershell
$version = (Invoke-RestMethod "https://api.github.com/repos/enolalabs/dotagen/releases/latest").tag_name
$url = "https://github.com/enolalabs/dotagen/releases/download/$version/dotagen_$($version.Substring(1))_windows_amd64.exe"
mkdir "$env:USERPROFILE\bin" -Force -ErrorAction SilentlyContinue
Invoke-WebRequest $url -OutFile "$env:USERPROFILE\bin\dotagen.exe"
```

Then add `%USERPROFILE%\bin` to your PATH if not already there.

### Install via `go install`

If you have Go 1.26+:

```bash
go install github.com/enolalabs/dotagen/cmd/dotagen@latest
```

### Build from source

<details>
<summary>Manual build instructions</summary>

#### Linux

```bash
git clone https://github.com/enolalabs/dotagen.git
cd dotagen
make build
sudo cp dotagen /usr/local/bin/
dotagen --version
```

#### macOS

```bash
git clone https://github.com/enolalabs/dotagen.git
cd dotagen
make build
sudo cp dotagen /usr/local/bin/
# Or: sudo cp dotagen /opt/homebrew/bin/
dotagen --version
```

#### Windows (PowerShell)

```powershell
git clone https://github.com/enolalabs/dotagen.git
cd dotagen
go build -ldflags="-X github.com/enolalabs/dotagen/v2/internal/cli/version=dev" -o dotagen.exe .\cmd\dotagen
mkdir "$env:USERPROFILE\bin" -Force -ErrorAction SilentlyContinue
Copy-Item dotagen.exe "$env:USERPROFILE\bin\"
dotagen --version
```

</details>

## Built-in Skills

dotagen ships with **841 skills** from **58 vendors** and **203 built-in agents**, sourced from the [awesome-agent-skills](https://github.com/enolalabs/awesome-agent-skills) registry, [obra/superpowers](https://github.com/obra/superpowers), [VoltAgent](https://github.com/VoltAgent/awesome-claude-code-subagents) and [claude-code-game-studios](https://github.com/donchitos/claude-code-game-studios). They are injected automatically when you run `dotagen init`.

All skills are **disabled by default**. You decide which skills to enable and for which platforms.

### Categories

| Category | Count |
|---|---|
| Developer Tools | 191 |
| Product & Strategy | 137 |
| AI & Machine Learning | 87 |
| Game Development | 73 |
| Testing & QA | 60 |
| Cloud & Infrastructure | 56 |
| Backend & APIs | 53 |
| DevOps & Monitoring | 35 |
| Databases & Data | 24 |
| Frontend & UI | 23 |
| Productivity & Collaboration | 32 |
| Security | 21 |
| Documents & Content | 18 |
| Search & Web | 16 |
| Mobile & Desktop | 15 |

### Vendors

<details>
<summary><strong>View all 58 vendors</strong></summary>

| Vendor | Skills | Category |
|---|---|---|
| Microsoft | 132 | Developer Tools |
| Game Studios | 73 | Game Development |
| phuryn | 65 | Product & Strategy |
| testmu-ai | 48 | Testing & QA |
| deanpeters | 46 | Product & Strategy |
| OpenAI | 36 | AI & Machine Learning |
| Sentry | 28 | DevOps & Monitoring |
| Superpowers | 14 | Developer Tools |
| Matt Pocock | 26 | Developer Tools |
| Garry Tan | 27 | Developer Tools |
| Trail of Bits | 21 | Security |
| Google | 19 | Cloud & Infrastructure |
| Venice AI | 19 | AI & Machine Learning |
| Anthropic | 17 | Documents & Content |
| Google Workspace | 17 | Productivity & Collaboration |
| Auth0 | 14 | Backend & APIs |
| Corey Haines | 14 | Product & Strategy |
| Apollo GraphQL | 13 | Backend & APIs |
| WordPress | 13 | Mobile & Desktop |
| Kim Barrett | 12 | Product & Strategy |
| Brave | 11 | Search & Web |
| HashiCorp | 11 | Cloud & Infrastructure |
| Netlify | 11 | Cloud & Infrastructure |
| NVIDIA | 11 | AI & Machine Learning |
| MiniMax AI | 10 | AI & Machine Learning |
| Cloudflare | 8 | Cloud & Infrastructure |
| GreenSock | 8 | Frontend & UI |
| Firebase | 7 | Cloud & Infrastructure |
| Hugging Face | 7 | AI & Machine Learning |
| Addy Osmani | 6 | Frontend & UI |
| Binance | 6 | Backend & APIs |
| Browserbase | 6 | Testing & QA |
| DuckDB | 6 | Databases & Data |
| MongoDB | 6 | Databases & Data |
| Better Auth | 5 | Backend & APIs |
| ClickHouse | 5 | Databases & Data |
| Datadog | 5 | DevOps & Monitoring |
| Firecrawl | 5 | Search & Web |
| Resend | 5 | Backend & APIs |
| Figma | 4 | Frontend & UI |
| Notion | 4 | Productivity & Collaboration |
| Sanity | 4 | Backend & APIs |
| Tinybird | 4 | Databases & Data |
| VoltAgent | 4 | Developer Tools |
| Callstack | 3 | Frontend & UI |
| Cypress | 3 | Testing & QA |
| Google Gemini | 3 | AI & Machine Learning |
| Neon | 3 | Databases & Data |
| Angular | 2 | Frontend & UI |
| CodeRabbit | 2 | DevOps & Monitoring |
| Expo | 2 | Mobile & Desktop |
| Stripe | 2 | Backend & APIs |
| Zero | 2 | Backend & APIs |
| fal.ai | 1 | AI & Machine Learning |
| Remotion | 1 | Documents & Content |
| Courier | 1 | Backend & APIs |
| Typefully | 1 | Backend & APIs |

</details>

### Skill Naming Convention

Built-in skills use the `dotagent:` prefix:

- **Directory:** `dotagent-{vendor}-{skill-name}/`
- **Frontmatter name:** `dotagent:{vendor}:{skill-name}`
- **Example:** `dotagent-anthropics-pdf/` → `dotagent:anthropics:pdf`

## Quick Start

### 1. Initialize

```bash
dotagen init
```

Creates `.dotagen/` with all 841 built-in skills and a config file where everything is disabled by default:

```
.dotagen/
├── config.yaml       # Configuration — set targets to enable skills
├── agents/           # Your custom agent definitions (*.md)
├── skills/           # 841 built-in skill directories (dotagent-*/SKILL.md)
├── .generated/       # Rendered output (git-ignored)
└── .gitignore
```

### 2. Configure

Edit `.dotagen/config.yaml` to enable skills by setting their `targets`. By default all skills have `targets: []` (disabled):

```yaml
targets:
  - claude-code
  - cursor
  - gemini-cli
  - opencode

agents:
  # Add your own agents here

skills:
  dotagent-anthropics-pdf:
    targets: all                    # Enable for all platforms
  dotagent-openai-doc:
    targets: all
  dotagent-microsoft-csharp:
    targets:
      - claude-code                 # Claude Code only
  # ... all other skills remain disabled (targets: [])
```

**`targets` values:**
- `all` — Apply to all configured platforms
- A specific list, e.g. `[claude-code, gemini-cli]`
- `[]` — Disabled (not synced)

### 3. Sync

```bash
# Sync all platforms
dotagen sync

# Sync a single platform
dotagen sync claude-code
```

Output:

```
✓ Synced 3 skills to 3 platform(s)

  claude-code:
    ✓ .claude/skills/dotagent-anthropics-pdf/SKILL.md → .dotagen/.generated/claude-code/...
    ✓ .claude/skills/dotagent-openai-doc/SKILL.md → .dotagen/.generated/claude-code/...

  gemini-cli:
    ✓ .agents/skills/dotagent-anthropics-pdf/SKILL.md → .dotagen/.generated/gemini-cli/...
    ✓ .agents/skills/dotagent-openai-doc/SKILL.md → .dotagen/.generated/gemini-cli/...

  ...
```

### 4. Check Status

```bash
dotagen status
```

Shows the state of each agent/skill on each platform:

- `✓ synced` — Up to date
- `⚠ out-of-date` — Source has changed, needs re-sync
- `✗ missing` — Not yet created
- `💔 broken` — Symlink is broken

## CLI Commands

| Command | Description |
|---|---|
| `dotagen init` | Initialize `.dotagen/` with 841 built-in skills (all disabled) |
| `dotagen sync [target]` | Render & symlink agents and skills. Optionally specify a target platform |
| `dotagen status` | Show sync status of all agents and skills |
| `dotagen clean` | Remove all generated files and symlinks (agents + skills) |
| `dotagen skill list` | List all skills with categories and targets |
| `dotagen skill create <name>` | Create a new skill directory with scaffold SKILL.md |
| `dotagen skill delete <name>` | Delete a skill directory and config entry |
| `dotagen serve` | Start web dashboard at `http://localhost:7890` |
| `dotagen update` | Update dotagen to the latest version |
| `dotagen version` | Print version |

## Web Dashboard

```bash
dotagen serve
```

Starts a web dashboard at `http://localhost:7890` with the following features:

- **Overview** — Stats, category distribution, platform health, symlink status
- **Agent Library** — Matrix view: toggle agents across 8 platforms, bulk column toggle
- **Skill Library** — Matrix view: toggle 768 skills across 8 platforms, filter by category/vendor
- **Status** — View all symlinks (agents + skills) with type badges, Fix Broken button
- **Preview** — View rendered output for each platform
- **Clean / Reinit** — Remove all symlinks or reinitialize from built-in library

### Screenshots

<details>
<summary><strong>Overview Dashboard</strong></summary>

![Overview](docs/screenshots/overview.png)

</details>

<details>
<summary><strong>Agent Library (Matrix View)</strong></summary>

![Agent Library](docs/screenshots/agent-library.png)

</details>

<details>
<summary><strong>Skill Library (Matrix View)</strong></summary>

![Skill Library](docs/screenshots/skill-library.png)

</details>

<details>
<summary><strong>Symlink Status</strong></summary>

![Status](docs/screenshots/status.png)

</details>

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
| `GET` | `/api/skills` | List all skills |
| `GET` | `/api/skills/{name}` | Get skill detail |
| `POST` | `/api/skills` | Create a new skill |
| `PUT` | `/api/skills/{name}` | Update a skill |
| `DELETE` | `/api/skills/{name}` | Delete a skill |
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
    disabled: true | false   # (optional) Temporarily disable

# Skill-to-target mapping
skills:
  <skill-name>:
    targets: all | [target1, target2, ...]
    disabled: true | false   # (optional) Temporarily disable
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

**Note:** Frontmatter (`---`) is optional. When present, the `description` field is used by some platform adapters when rendering output.

### Skill Format

Skills are directory-based with a `SKILL.md` entry point:

```
skills/
└── my-skill/
    ├── SKILL.md           # Required — skill definition
    ├── references/        # Optional — reference files
    │   └── api-docs.md
    └── examples/          # Optional — example files
        └── demo.md
```

```markdown
---
name: my-skill
description: What this skill does
category: engineering
---

# My Skill

Instructions for the skill...
```

## Architecture

```
.dotagen/agents/*.md  +  .dotagen/skills/dotagent-*/  +  .dotagen/config.yaml
          │                      │
          ▼                      ▼
     Config Parser → Agent Parser / Skill Parser → Renderer
                                                      │
                                            ┌─────────┴─────────┐
                                            ▼                   ▼
                               .generated/{platform}/   Skill dir symlinks
                                            │                   │
                                            ▼                   ▼
                               Symlink → .claude/agents/    .claude/skills/
                                        .agents/            .agents/skills/
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
│   ├── skill/                   # Skill directory parser
│   ├── builtin/                 # Built-in agents (go:embed)
│   │   ├── embed.go
│   │   └── agents/              # Agent .md files
│   ├── cli/                     # CLI commands (cobra)
│   ├── config/                  # Config parser & validator
│   ├── engine/                  # Renderer & symlink manager
│   │   ├── renderer.go          # Agent renderer
│   │   └── skill_renderer.go    # Skill renderer
│   ├── platform/                # Platform adapters
│   │   ├── adapter.go           # Agent adapter interface
│   │   ├── skill_adapter.go     # Skill adapter interface
│   │   ├── antigravity.go
│   │   ├── claude_code.go
│   │   ├── codex.go
│   │   ├── cursor.go
│   │   ├── gemini_cli.go
│   │   ├── github_copilot.go
│   │   ├── opencode.go
│   │   ├── windsurf.go
│   │   └── registry.go
│   └── web/                     # Web dashboard (go:embed)
│       ├── server.go
│       ├── api.go               # Agent API
│       ├── skill_api.go         # Skill API
│       └── static/
├── skillsrc/                    # Built-in skills (go:embed)
│   ├── embed.go
│   └── data/                    # 841 skill directories
├── scripts/
│   └── fetch-official-skills.py # Skill fetcher from awesome-agent-skills
├── go.mod
├── Makefile
└── README.md
```

## Acknowledgments

- The built-in skills are sourced from the [**awesome-agent-skills**](https://github.com/enolalabs/awesome-agent-skills) registry — a curated collection of official vendor skills from 55 organizations including Microsoft, OpenAI, Anthropic, Google, NVIDIA, Stripe, Cloudflare, and many more — plus [**obra/superpowers**](https://github.com/obra/superpowers) and [**mattpocock/skills**](https://github.com/mattpocock/skills).
- The 154 `dotagent-voltagent-*` agents come from [**VoltAgent/awesome-claude-code-subagents**](https://github.com/VoltAgent/awesome-claude-code-subagents).
- The 49 `dotagent-game-studios-*` agents and 73 `dotagent:game-studios:*` skills come from [**donchitos/claude-code-game-studios**](https://github.com/donchitos/claude-code-game-studios) (MIT) — a full virtual game studio (Unity, Unreal, Godot specialists, designers, producers, QA, live-ops) with sprint/gate/release workflows.

## License

MIT
