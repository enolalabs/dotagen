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
git clone https://github.com/enolalabs/dotagen.git
cd dotagen
make build
```

The binary will be created at `./dotagen`.

### Or install directly

```bash
go install github.com/enolalabs/dotagen/cmd/dotagen@latest
```

## Built-in Agents

dotagen ships with **144 built-in agents** covering a wide range of specialties. They are injected automatically when you run `dotagen init`, so you don't need to create agents from scratch.

All agents are **disabled by default**. You decide which agents to enable and for which platforms.

<details>
<summary><strong>Business & Product</strong> (12 agents)</summary>

| Agent | Description |
|---|---|
| `business-analyst` | Analyze business processes, gather requirements, identify improvement opportunities |
| `content-marketer` | Content strategies, SEO-optimized marketing, multi-channel campaigns |
| `customer-success-manager` | Customer health assessment, retention strategies, upsell opportunities |
| `legal-advisor` | Draft contracts, review compliance, IP protection, legal risk assessment |
| `license-engineer` | OSI standard selection, dependency compliance, proprietary deployment |
| `marketing-analyst` | Campaign performance, attribution models, growth strategies |
| `product-manager` | Feature prioritization, roadmap planning, stakeholder alignment |
| `sales-engineer` | Technical pre-sales, solution architecture, proof-of-concept |
| `scrum-master` | Sprint planning, retrospectives, impediment removal, velocity improvement |
| `technical-writer` | API references, user guides, SDK documentation |
| `ux-researcher` | User research, usability testing, persona development |
| `wordpress-master` | WordPress architecture, WooCommerce, performance, security hardening |

</details>

<details>
<summary><strong>Core Development</strong> (11 agents)</summary>

| Agent | Description |
|---|---|
| `api-designer` | API specifications, RESTful patterns, GraphQL schema design |
| `backend-developer` | Server-side APIs, microservices, robust backend systems |
| `database-architect` | Schema design, query optimization, migration strategies |
| `frontend-developer` | Modern frontend with React/Vue/Angular, responsive design |
| `fullstack-developer` | End-to-end application development |
| `graphql-developer` | GraphQL schemas, resolvers, federation, subscriptions |
| `legacy-modernizer` | Modernize legacy systems, migration strategies |
| `low-level-designer` | OOP/functional class-level design, SOLID principles |
| `microservices-architect` | Distributed systems, service mesh, event-driven architecture |
| `ui-designer` | Visual interfaces, design systems, component libraries |
| `websocket-engineer` | Real-time bidirectional communication at scale |

</details>

<details>
<summary><strong>Data & AI</strong> (13 agents)</summary>

| Agent | Description |
|---|---|
| `ai-engineer` | End-to-end AI systems, model selection, deployment pipelines |
| `computer-vision-engineer` | Image/video analysis, object detection, OCR |
| `data-engineer` | ETL pipelines, data warehousing, streaming architecture |
| `data-pipeline-architect` | Large-scale data infrastructure, real-time processing |
| `data-scientist` | Statistical modeling, ML experiments, data visualization |
| `data-visualization` | Interactive dashboards, D3.js, Plotly, chart design |
| `elasticsearch-specialist` | Search clusters, query optimization, index management |
| `etl-specialist` | Data extraction, transformation, loading pipelines |
| `llm-architect` | LLM-powered applications, RAG, fine-tuning, prompt engineering |
| `ml-engineer` | ML model development, training pipelines, deployment |
| `nlp-engineer` | Text processing, sentiment analysis, language models |
| `playwright-expert` | Browser automation, E2E testing, scraping with Playwright |
| `prompt-engineer` | Prompt design, chain-of-thought, evaluation frameworks |

</details>

<details>
<summary><strong>Developer Experience</strong> (14 agents)</summary>

| Agent | Description |
|---|---|
| `build-engineer` | Build performance, compilation optimization, scaling |
| `cli-developer` | Command-line tools and terminal applications |
| `documentation-engineer` | Documentation-as-code, API docs, architecture docs |
| `git-specialist` | Advanced Git workflows, branching strategies, history management |
| `github-actions-specialist` | CI/CD with GitHub Actions, workflow optimization |
| `ide-plugin-developer` | IDE extension development for VS Code, JetBrains |
| `json-wrangler` | JSON/YAML transformation, schema validation, jq expert |
| `monorepo-engineer` | Nx/Turborepo/Lerna monorepo architecture |
| `open-source-advisor` | OSS contribution, governance, community building |
| `refactoring-specialist` | Code refactoring, tech debt reduction, pattern migration |
| `regex-master` | Complex regex patterns, validation, text extraction |
| `slack-expert` | Slack applications, bot development, API integrations |
| `tooling-engineer` | Developer tools, CLIs, code generators, build tools |
| `vibe-coder` | Rapid prototyping, creative coding, quick iteration |

</details>

<details>
<summary><strong>Infrastructure</strong> (16 agents)</summary>

| Agent | Description |
|---|---|
| `azure-infra-engineer` | Azure infrastructure, networking, deployment |
| `cicd-engineer` | CI/CD pipeline design, deployment automation |
| `cloud-architect` | Cloud infrastructure design, multi-cloud strategy |
| `devops-engineer` | Infrastructure automation, monitoring, deployment |
| `docker-expert` | Container optimization, multi-stage builds, Compose |
| `gcp-specialist` | Google Cloud Platform services and architecture |
| `kubernetes-specialist` | K8s cluster management, Helm, operators |
| `linux-sysadmin` | Linux server administration, shell scripting |
| `network-engineer` | Network architecture, security, troubleshooting |
| `nginx-specialist` | Nginx configuration, load balancing, reverse proxy |
| `powershell-admin` | Windows automation, Active Directory, system management |
| `security-engineer` | Security solutions, zero-trust architecture, CI/CD security |
| `sre-engineer` | SLO/SLI frameworks, error budgets, chaos engineering |
| `terraform-engineer` | Infrastructure as code, multi-cloud Terraform |
| `terragrunt-expert` | Terragrunt orchestration, DRY configurations |
| `windows-infra-admin` | Windows Server, Active Directory, Group Policy |

</details>

<details>
<summary><strong>Language Specialists</strong> (30 agents)</summary>

| Agent | Description |
|---|---|
| `angular-architect` | Angular 15+ enterprise applications |
| `astro-developer` | Astro framework, content-driven websites |
| `cpp-systems-developer` | C++ systems programming, memory management |
| `csharp-dotnet-developer` | C#/.NET enterprise applications |
| `django-developer` | Django web applications and REST APIs |
| `elixir-phoenix-developer` | Elixir/Phoenix real-time applications |
| `flutter-developer` | Flutter cross-platform mobile/web apps |
| `golang-pro` | Go applications, concurrency, performance |
| `java-enterprise-architect` | Java enterprise systems, Spring, microservices |
| `kotlin-expert` | Kotlin/Android development, coroutines |
| `laravel-expert` | Laravel PHP applications, Eloquent ORM |
| `nestjs-architect` | NestJS enterprise backends, microservices |
| `nextjs-developer` | Next.js full-stack applications, SSR/SSG |
| `nuxt-specialist` | Nuxt 3 applications, Vue ecosystem |
| `perl-modernizer` | Perl modernization, Moose, async patterns |
| `php-engineer` | Modern PHP 8+, frameworks, performance |
| `python-pro` | Python applications, async, data processing |
| `r-statistician` | R statistical computing, data analysis |
| `rails-developer` | Ruby on Rails applications |
| `react-native-developer` | React Native cross-platform mobile apps |
| `react-specialist` | React 18+, hooks, state management |
| `ruby-pro` | Ruby applications, metaprogramming |
| `rust-engineer` | Rust systems programming, memory safety |
| `spring-boot-engineer` | Spring Boot 3+ enterprise applications |
| `sql-pro` | SQL optimization, schema design, indexing |
| `swift-expert` | Swift/iOS/macOS native applications |
| `symfony-specialist` | Symfony 6+/7+ applications, Doctrine ORM |
| `typescript-pro` | TypeScript advanced type patterns |
| `vue-expert` | Vue 3 Composition API, Nuxt development |
| `wordpress-master` | WordPress themes, plugins, WooCommerce |

</details>

<details>
<summary><strong>Meta-Orchestration</strong> (11 agents)</summary>

| Agent | Description |
|---|---|
| `agent-installer` | Discover, browse, install Claude Code agents |
| `agent-organizer` | Assemble and optimize multi-agent teams |
| `codebase-orchestrator` | Repository-wide refactor governance |
| `context-manager` | Shared state management across agents |
| `error-coordinator` | Coordinated error handling across components |
| `it-ops-orchestrator` | Multi-domain IT operations orchestration |
| `knowledge-synthesizer` | Extract patterns from agent interactions |
| `multi-agent-coordinator` | Coordinate concurrent agents |
| `performance-monitor` | Observability infrastructure, metrics tracking |
| `task-distributor` | Task distribution, queue management, load balancing |
| `workflow-orchestrator` | Business process workflow design |

</details>

<details>
<summary><strong>Quality & Security</strong> (16 agents)</summary>

| Agent | Description |
|---|---|
| `accessibility-tester` | WCAG compliance, accessibility testing |
| `ad-security-reviewer` | Active Directory security posture audit |
| `ai-writing-auditor` | Audit and rewrite AI-generated content |
| `architect-reviewer` | System design review, architectural patterns |
| `chaos-engineer` | Controlled failure experiments, resilience |
| `code-reviewer` | Comprehensive code review, security, quality |
| `compliance-auditor` | Regulatory compliance, audit controls |
| `debugger` | Bug diagnosis, root cause analysis |
| `error-detective` | Error correlation, failure chain analysis |
| `penetration-tester` | Authorized security penetration testing |
| `performance-engineer` | Performance bottleneck identification |
| `powershell-security-hardening` | PowerShell security, remoting hardening |
| `qa-expert` | QA strategy, test planning, coverage |
| `security-auditor` | Security audits, compliance assessments |
| `test-automator` | Automated test frameworks, CI/CD testing |
| `ui-ux-tester` | UI/UX functionality testing |

</details>

<details>
<summary><strong>Research & Analysis</strong> (8 agents)</summary>

| Agent | Description |
|---|---|
| `competitive-analyst` | Competitor analysis, market benchmarking |
| `data-researcher` | Multi-source data discovery and validation |
| `market-researcher` | Market analysis, consumer behavior |
| `project-idea-validator` | Idea pressure-testing, competitor teardown |
| `research-analyst` | Multi-source research synthesis |
| `scientific-literature-researcher` | Scientific literature search, structured data |
| `search-specialist` | Advanced search strategies, query optimization |
| `trend-analyst` | Emerging patterns, industry shift prediction |

</details>

<details>
<summary><strong>Specialized Domains</strong> (13 agents)</summary>

| Agent | Description |
|---|---|
| `api-documenter` | API documentation, OpenAPI specifications |
| `blockchain-developer` | Smart contracts, DApps, blockchain protocols |
| `embedded-systems` | Firmware, RTOS, microcontroller development |
| `fintech-engineer` | Payment systems, financial compliance |
| `game-developer` | Game systems, graphics, multiplayer networking |
| `healthcare-admin` | Healthcare administration, HIPAA compliance |
| `iot-engineer` | IoT device management, edge computing |
| `m365-admin` | Microsoft 365 administration automation |
| `mobile-app-developer` | iOS/Android mobile application development |
| `payment-integration` | Payment gateway integration, PCI compliance |
| `quant-analyst` | Quantitative trading, financial modeling |
| `risk-manager` | Enterprise risk identification and mitigation |
| `seo-specialist` | SEO audits, keyword strategy, optimization |

</details>

## Built-in Skills

dotagen also ships with **16 built-in skills** (slash commands) sourced from [mattpocock/skills](https://github.com/mattpocock/skills). Skills are directory-based (`ds-<name>/SKILL.md`) and are injected alongside agents during `dotagen init`.

| Skill | Category | Description |
|---|---|---|
| `ds-caveman` | productivity | Simplify explanations to the most basic level |
| `ds-diagnose` | engineering | Disciplined diagnosis loop: reproduce → minimise → hypothesise → instrument → fix → regression-test |
| `ds-git-guardrails` | misc | Set up hooks to block dangerous git commands (push, reset --hard, clean, etc.) |
| `ds-grill-me` | productivity | Interview the user relentlessly about a plan or design until shared understanding |
| `ds-grill-with-docs` | engineering | Challenge your plan against the existing domain model and update documentation inline |
| `ds-improve-codebase-architecture` | engineering | Find deepening opportunities informed by domain language and ADRs |
| `ds-migrate-to-shoehorn` | misc | Migrate test files from `as` assertions to @total-typescript/shoehorn |
| `ds-scaffold-exercises` | misc | Create exercise directory structures with sections, problems, solutions |
| `ds-setup-matt-pocock-skills` | engineering | Set up AGENTS.md skill block and docs/agents/ for repo-specific context |
| `ds-setup-pre-commit` | misc | Set up Husky pre-commit hooks with lint-staged, type checking, tests |
| `ds-tdd` | engineering | Test-driven development with red-green-refactor loop |
| `ds-to-issues` | engineering | Break a plan/spec/PRD into independently-grabbable issues using tracer-bullet slices |
| `ds-to-prd` | engineering | Turn the current conversation context into a PRD |
| `ds-triage` | engineering | Triage issues through a state machine driven by triage roles |
| `ds-write-a-skill` | productivity | Create new agent skills with proper structure and progressive disclosure |
| `ds-zoom-out` | engineering | Zoom out for broader context and higher-level perspective on code |

## Quick Start

### 1. Initialize

```bash
dotagen init
```

Creates `.dotagen/` with all 144 built-in agents, 16 built-in skills, and a config file where everything is disabled by default:

```
.dotagen/
├── config.yaml       # Configuration — set targets to enable agents/skills
├── agents/           # 144 built-in agent definitions (*.md)
├── skills/           # 16 built-in skill directories (ds-*/SKILL.md)
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

skills:
  ds-diagnose:
    targets: all              # Enable for all platforms
  ds-tdd:
    targets:
      - claude-code
  # ... all other skills remain disabled (targets: [])
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
| `dotagen init` | Initialize `.dotagen/` with 144 built-in agents and 16 skills (all disabled) |
| `dotagen sync [target]` | Render & symlink agents and skills. Optionally specify a target platform |
| `dotagen status` | Show sync status of all agents and skills |
| `dotagen clean` | Remove all generated files and symlinks (agents + skills) |
| `dotagen skill list` | List all skills with categories and targets |
| `dotagen skill create <name>` | Create a new skill directory with scaffold SKILL.md |
| `dotagen skill delete <name>` | Delete a skill directory and config entry |
| `dotagen serve` | Start web dashboard at `http://localhost:7890` |
| `dotagen --version` | Print version |
| `dotagen --help` | Print help |

## Web Dashboard

```bash
dotagen serve
```

Starts a web dashboard at `http://localhost:7890` with the following features:

- **Agent Management** — CRUD agents directly from the UI
- **Skill Management** — CRUD skills with multi-category dropdown
- **Target Matrix** — Assign/unassign agents and skills to platforms via toggle
- **Preview** — View rendered output for each platform
- **Sync/Clean** — Trigger sync or clean from the web UI (applies to both agents and skills)
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

**Note:** Frontmatter (`---`) is optional. When present, the `description` field is used by the Cursor and OpenCode adapters when rendering output.

## Architecture

```
.dotagen/agents/*.md  +  .dotagen/skills/ds-*/  +  .dotagen/config.yaml
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
                                        .cursor/rules/     .cursor/skills/
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
│   ├── builtin/                 # Built-in agents & skills (go:embed)
│   │   ├── embed.go             # Agent embed
│   │   ├── skill_embed.go       # Skill embed
│   │   ├── agents/              # 144 agent .md files
│   │   └── skills/              # 16 skill directories
│   ├── cli/                     # CLI commands (cobra)
│   ├── config/                  # Config parser & validator
│   ├── engine/                  # Renderer & symlink manager
│   │   ├── renderer.go          # Agent renderer
│   │   └── skill_renderer.go    # Skill renderer
│   ├── platform/                # Platform adapters
│   │   ├── adapter.go           # Agent adapter interface
│   │   ├── skill_adapter.go     # Skill adapter interface
│   │   ├── claude_code.go
│   │   ├── cursor.go
│   │   ├── gemini_cli.go
│   │   ├── opencode.go
│   │   └── registry.go
│   └── web/                     # Web dashboard (go:embed)
│       ├── server.go
│       ├── api.go               # Agent API
│       ├── skill_api.go         # Skill API
│       └── static/
├── go.mod
├── Makefile
└── README.md
```

## Acknowledgments

- The 144 built-in agents are sourced from the [**VoltAgent**](https://github.com/VoltAgent/voltagent) project — an excellent collection of production-ready sub-agent prompts.
- The 16 built-in skills are sourced from [**mattpocock/skills**](https://github.com/mattpocock/skills) — a curated set of slash commands for engineering workflows.

## License

MIT
