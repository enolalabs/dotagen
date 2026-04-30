# Dotagen — Implementation Plan

> **Define sub-agents once, inject everywhere.**
>
> A Go CLI tool with embedded web dashboard that lets you define coding sub-agents
> in markdown and inject them into multiple coding agent platforms.

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Phase 1: Project Scaffold & Foundation](#phase-1-project-scaffold--foundation)
- [Phase 2: Core Engine](#phase-2-core-engine)
- [Phase 3: CLI Commands](#phase-3-cli-commands)
- [Phase 4: Platform Adapters](#phase-4-platform-adapters)
- [Phase 5: Web UI (Embedded)](#phase-5-web-ui-embedded)
- [Phase 6: Polish & Release](#phase-6-polish--release)
- [Appendix](#appendix)

---

## Overview

### Problem

Mỗi coding agent (Cursor, Claude Code, Gemini CLI, OpenCode...) có cách cấu hình sub-agents riêng.
Điều này buộc developer phải:
- Viết lại cùng một instructions ở nhiều format khác nhau
- Maintain N bộ config cho N agent
- Setup lại mỗi khi chuyển agent hoặc thêm agent mới

### Solution

**dotagen** — một CLI tool viết bằng Go, cho phép:
1. Define sub-agents **một lần** bằng markdown (`.dotagen/agents/*.md`)
2. Cấu hình wiring bằng YAML (`.dotagen/config.yaml`)
3. CLI **render + symlink** vào đúng vị trí cho từng coding agent platform
4. Web dashboard local để quản lý trực quan

### Tech Stack

| Component | Technology |
|---|---|
| Language | Go 1.26.1 |
| CLI framework | `cobra` |
| Config parsing | `gopkg.in/yaml.v3` |
| Markdown parsing | `github.com/yuin/goldmark` (nếu cần parse frontmatter) |
| Frontmatter | `github.com/adrg/frontmatter` |
| Web UI | Vanilla HTML/CSS/JS, embedded via `go:embed` |
| HTTP server | `net/http` (stdlib) |
| Testing | `testing` (stdlib) + `testify` |

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      dotagen CLI                         │
│                                                          │
│  ┌──────────┐  ┌──────────────┐  ┌───────────────────┐  │
│  │  Commands │  │  Web Server  │  │  Platform Adapters│  │
│  │  (cobra)  │  │  (go:embed)  │  │                   │  │
│  │           │  │              │  │  ┌─────────────┐  │  │
│  │  init     │  │  REST API    │  │  │ Claude Code │  │  │
│  │  sync     │  │  Static UI   │  │  │ Cursor      │  │  │
│  │  clean    │  │              │  │  │ Gemini CLI  │  │  │
│  │  status   │  │              │  │  │ OpenCode    │  │  │
│  │  serve    │  │              │  │  └─────────────┘  │  │
│  └─────┬─────┘  └──────┬──────┘  └────────┬──────────┘  │
│        │               │                   │             │
│  ┌─────▼───────────────▼───────────────────▼──────────┐  │
│  │                  Core Engine                        │  │
│  │                                                     │  │
│  │  ┌──────────┐ ┌──────────┐ ┌────────┐ ┌─────────┐  │  │
│  │  │  Config  │ │  Agent   │ │Template│ │ Symlink │  │  │
│  │  │  Parser  │ │  Parser  │ │Renderer│ │ Manager │  │  │
│  │  └──────────┘ └──────────┘ └────────┘ └─────────┘  │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                          │
└──────────────────────────┬───────────────────────────────┘
                           │
                           ▼
              ┌─────────────────────────┐
              │    .dotagen/ directory    │
              │                          │
              │  config.yaml             │
              │  agents/                 │
              │    review-code.md        │
              │    planning.md           │
              │    testing.md            │
              │  .generated/             │
              │    claude-code/           │
              │    cursor/               │
              │    gemini-cli/           │
              │    opencode/             │
              └─────────────────────────┘
```

### Data Flow

```
[.dotagen/agents/*.md] + [.dotagen/config.yaml]
        │
        ▼
   Config Parser → load targets, agent-target mapping
        │
        ▼
   Agent Parser → parse markdown + optional frontmatter
        │
        ▼
   Template Renderer → dịch sang format của từng platform
        │
        ▼
   .dotagen/.generated/{platform}/{agent}.{ext}
        │
        ▼
   Symlink Manager → tạo symlink từ platform path → generated file
        │
        ▼
   .claude/agents/review-code.md → symlink → .dotagen/.generated/claude-code/review-code.md
   .cursor/rules/review-code.mdc → symlink → .dotagen/.generated/cursor/review-code.mdc
   ...
```

---

## Project Structure

```
dotagen/
├── cmd/
│   └── dotagen/
│       └── main.go                  # Entry point
│
├── internal/
│   ├── config/
│   │   ├── config.go                # Config types & parsing
│   │   └── config_test.go
│   │
│   ├── agent/
│   │   ├── parser.go                # Agent markdown parser
│   │   ├── parser_test.go
│   │   └── types.go                 # Agent types
│   │
│   ├── engine/
│   │   ├── renderer.go              # Template rendering engine
│   │   ├── renderer_test.go
│   │   ├── symlink.go               # Symlink creation/management
│   │   └── symlink_test.go
│   │
│   ├── platform/
│   │   ├── adapter.go               # Platform adapter interface
│   │   ├── claude_code.go           # Claude Code adapter
│   │   ├── claude_code_test.go
│   │   ├── cursor.go                # Cursor adapter
│   │   ├── cursor_test.go
│   │   ├── gemini_cli.go            # Gemini CLI adapter
│   │   ├── gemini_cli_test.go
│   │   ├── opencode.go              # OpenCode adapter
│   │   ├── opencode_test.go
│   │   └── registry.go              # Platform registry
│   │
│   ├── cli/
│   │   ├── root.go                  # Root cobra command
│   │   ├── init.go                  # dotagen init
│   │   ├── sync.go                  # dotagen sync [target]
│   │   ├── clean.go                 # dotagen clean
│   │   ├── status.go                # dotagen status
│   │   └── serve.go                 # dotagen serve
│   │
│   └── web/
│       ├── server.go                # HTTP server + API handlers
│       ├── api.go                   # REST API endpoints
│       └── static/                  # Embedded web assets
│           ├── index.html
│           ├── style.css
│           └── app.js
│
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── IMPLEMENTATION_PLAN.md           # ← This file
```

---

## Phase 1: Project Scaffold & Foundation

> **Goal**: Go module sẵn sàng, build được, có base CLI structure.

### Tasks

#### 1.1 — Initialize Go module

```bash
go mod init github.com/k0walski/dotagen
```

#### 1.2 — Install dependencies

```bash
go get github.com/spf13/cobra@latest
go get gopkg.in/yaml.v3
go get github.com/adrg/frontmatter
go get github.com/stretchr/testify
```

#### 1.3 — Create entry point

- `cmd/dotagen/main.go` — khởi tạo cobra root command, gọi `Execute()`

#### 1.4 — Create root CLI command

- `internal/cli/root.go` — root command với version flag, help text
- Banner/logo khi chạy `dotagen` không có subcommand

#### 1.5 — Define core types

- `internal/config/config.go`:

```go
type Config struct {
    Targets []string              `yaml:"targets"`
    Agents  map[string]AgentConfig `yaml:"agents"`
}

type AgentConfig struct {
    Targets  []string `yaml:"targets"`   // ["claude-code", "cursor"] hoặc ["all"]
    Disabled bool     `yaml:"disabled"`  // optional: tạm disable agent
}
```

- `internal/agent/types.go`:

```go
type Agent struct {
    Name        string            // Tên file (không có .md)
    Content     string            // Nội dung markdown (body)
    Frontmatter map[string]string // Optional frontmatter metadata
    FilePath    string            // Đường dẫn gốc
}
```

### Acceptance Criteria

- [ ] `go build ./cmd/dotagen` thành công
- [ ] `./dotagen --version` in ra version
- [ ] `./dotagen --help` in ra usage

---

## Phase 2: Core Engine

> **Goal**: Đọc được config + agents, render được output, quản lý symlinks.

### Tasks

#### 2.1 — Config Parser (`internal/config/`)

- Load & validate `.dotagen/config.yaml`
- Resolve `targets: all` → expand thành danh sách tất cả targets
- Validate target names (chỉ chấp nhận: `claude-code`, `cursor`, `gemini-cli`, `opencode`)
- Auto-discover agents từ `.dotagen/agents/` directory

```go
func LoadConfig(dotgenDir string) (*Config, error)
func (c *Config) Validate() error
func (c *Config) ResolveTargets(agentName string) []string
```

#### 2.2 — Agent Parser (`internal/agent/`)

- Scan `.dotagen/agents/` directory
- Parse markdown files, tách frontmatter (nếu có) và body content
- Return `[]Agent`

```go
func ParseAgentsDir(agentsDir string) ([]Agent, error)
func ParseAgentFile(filePath string) (*Agent, error)
```

#### 2.3 — Template Renderer (`internal/engine/renderer.go`)

- Nhận `Agent` + target platform → render ra string output đúng format
- Ghi output vào `.dotagen/.generated/{platform}/{filename}`

```go
type Renderer struct {
    adapters map[string]platform.Adapter
}

func (r *Renderer) Render(agent Agent, target string) (string, error)
func (r *Renderer) RenderAll(agents []Agent, config *Config) error
```

#### 2.4 — Symlink Manager (`internal/engine/symlink.go`)

- Tạo symlink từ platform expected path → `.dotagen/.generated/...`
- Xóa symlink khi clean
- Detect broken symlinks
- Handle overwrite (remove existing → create new)

```go
func CreateSymlink(src, dst string) error
func RemoveSymlink(path string) error
func IsSymlink(path string) (bool, error)
func ListManagedSymlinks(projectDir string) ([]SymlinkInfo, error)
```

### Acceptance Criteria

- [ ] Parse `.dotagen/config.yaml` thành công với validation
- [ ] Parse `.dotagen/agents/*.md` với frontmatter support
- [ ] Render agent content sang đúng format cho mỗi platform
- [ ] Tạo/xóa symlinks chính xác
- [ ] Unit tests pass cho tất cả components

---

## Phase 3: CLI Commands

> **Goal**: Tất cả CLI commands hoạt động end-to-end.

### Tasks

#### 3.1 — `dotagen init` (`internal/cli/init.go`)

Tạo scaffold `.dotagen/` directory:

```
.dotagen/
  config.yaml           ← Template mặc định
  agents/
    example.md          ← Agent mẫu để user tham khảo
  .generated/           ← Empty, git-ignored
  .gitignore            ← Ignore .generated/
```

**Config mặc định:**

```yaml
# dotagen configuration
# Docs: https://github.com/k0walski/dotagen

targets:
  - claude-code
  - cursor
  - gemini-cli
  - opencode

agents:
  example:
    targets: all
```

**Agent mẫu (`agents/example.md`):**

```markdown
# Example Agent

You are a helpful coding assistant. This is an example agent
created by `dotagen init`. Edit this file or create new agents
in the `.dotagen/agents/` directory.

## Guidelines

- Follow project coding standards
- Write clean, maintainable code
- Add tests for new features
```

**Logic:**
- Check nếu `.dotagen/` đã tồn tại → hỏi có muốn overwrite không
- Tạo structure
- In ra hướng dẫn next steps

#### 3.2 — `dotagen sync [target]` (`internal/cli/sync.go`)

**Logic:**
1. Load config
2. Parse tất cả agents
3. Nếu có `[target]` arg → chỉ sync cho target đó
4. Nếu không có arg → sync tất cả targets
5. Render output cho mỗi (agent, target) pair
6. Tạo parent directories nếu chưa có
7. Tạo symlinks (overwrite nếu đã tồn tại)
8. In summary: bao nhiêu agents synced, bao nhiêu symlinks created

**Output mẫu:**

```
✓ Synced 3 agents to 4 platforms

  claude-code:
    ✓ .claude/agents/review-code.md → .dotagen/.generated/claude-code/review-code.md
    ✓ .claude/agents/testing.md → .dotagen/.generated/claude-code/testing.md

  cursor:
    ✓ .cursor/rules/review-code.mdc → .dotagen/.generated/cursor/review-code.mdc
    ✓ .cursor/rules/planning.mdc → .dotagen/.generated/cursor/planning.mdc
    ✓ .cursor/rules/testing.mdc → .dotagen/.generated/cursor/testing.mdc

  gemini-cli:
    ✓ .gemini/agents/testing.md → .dotagen/.generated/gemini-cli/testing.md

  opencode:
    ✓ .opencode/agents/planning.md → .dotagen/.generated/opencode/planning.md
    ✓ .opencode/agents/testing.md → .dotagen/.generated/opencode/testing.md
```

#### 3.3 — `dotagen clean` (`internal/cli/clean.go`)

**Logic:**
1. Load config
2. Tìm tất cả symlinks đã được tạo bởi dotagen
3. Remove symlinks
4. Remove `.dotagen/.generated/` contents
5. In summary

#### 3.4 — `dotagen status` (`internal/cli/status.go`)

**Logic:**
1. Load config
2. Check từng (agent, target) pair
3. Report: synced ✓, out-of-date ⚠, missing ✗, broken symlink 💔

**Output mẫu:**

```
dotagen status

Agents: 3 defined
Targets: 4 configured

  review-code:
    ✓ claude-code  (synced)
    ✓ cursor       (synced)
    ✗ gemini-cli   (not targeted)
    ✗ opencode     (not targeted)

  planning:
    ✗ claude-code  (not targeted)
    ⚠ cursor       (out of date)
    ✗ gemini-cli   (not targeted)
    ✓ opencode     (synced)

  testing:
    ✓ claude-code  (synced)
    ✓ cursor       (synced)
    ✓ gemini-cli   (synced)
    ✓ opencode     (synced)
```

**Out-of-date detection:** So sánh hash (MD5/SHA256) của generated file vs re-rendered content.

#### 3.5 — `dotagen serve` (`internal/cli/serve.go`)

- Khởi động HTTP server
- Serve embedded static files
- Expose REST API
- Default port: `7890`, configurable via `--port`
- Auto-open browser (optional, `--no-open`)

### Acceptance Criteria

- [ ] `dotagen init` tạo scaffold hoàn chỉnh
- [ ] `dotagen sync` render + symlink tất cả agents
- [ ] `dotagen sync cursor` chỉ sync Cursor
- [ ] `dotagen clean` xóa sạch generated files + symlinks
- [ ] `dotagen status` hiển thị trạng thái chính xác
- [ ] `dotagen serve` khởi động web server

---

## Phase 4: Platform Adapters

> **Goal**: Mỗi platform có adapter riêng, biết cách transform agent content sang đúng format.

### Interface

```go
// internal/platform/adapter.go
type Adapter interface {
    // Name returns the platform identifier (e.g., "claude-code")
    Name() string

    // Render transforms agent content into platform-specific format
    Render(agent agent.Agent) (string, error)

    // OutputPath returns where the generated file should be placed
    // relative to .dotagen/.generated/
    OutputPath(agentName string) string

    // SymlinkPath returns where the symlink should point to
    // relative to project root
    SymlinkPath(agentName string) string

    // EnsureDirectories creates necessary parent directories
    EnsureDirectories(projectDir string) error
}
```

### 4.1 — Claude Code Adapter

| Field | Value |
|---|---|
| Symlink path | `.claude/agents/{name}.md` |
| Output format | Pure markdown |
| Transform | Không cần transform, giữ nguyên content |

#### 4.2 — Cursor Adapter

| Field | Value |
|---|---|
| Symlink path | `.cursor/rules/{name}.mdc` |
| Output format | YAML frontmatter + markdown |
| Transform | Thêm frontmatter với `description`, `globs`, `alwaysApply` |

**Output mẫu:**

```markdown
---
description: Review code agent
globs:
alwaysApply: true
---

# Review Code

You are a senior code reviewer...
```

> **Note**: Cursor dùng `.mdc` extension (markdown component). Frontmatter fields:
> - `description`: Lấy từ dòng đầu tiên hoặc frontmatter của agent source
> - `alwaysApply: true`: Mặc định always apply
> - `globs`: Để trống (user có thể override trong config)

#### 4.3 — Gemini CLI Adapter

| Field | Value |
|---|---|
| Symlink path | `.gemini/agents/{name}.md` |
| Output format | Pure markdown |
| Transform | Không cần transform, giữ nguyên content |

#### 4.4 — OpenCode Adapter

| Field | Value |
|---|---|
| Symlink path | `.opencode/agents/{name}.md` |
| Output format | YAML frontmatter + markdown |
| Transform | Thêm frontmatter với `description`, `mode: subagent` |

**Output mẫu:**

```markdown
---
description: "Review code agent"
mode: subagent
---

# Review Code

You are a senior code reviewer...
```

#### 4.5 — Platform Registry

```go
// internal/platform/registry.go
type Registry struct {
    adapters map[string]Adapter
}

func NewRegistry() *Registry {
    r := &Registry{adapters: make(map[string]Adapter)}
    r.Register(NewClaudeCodeAdapter())
    r.Register(NewCursorAdapter())
    r.Register(NewGeminiCLIAdapter())
    r.Register(NewOpenCodeAdapter())
    return r
}

func (r *Registry) Get(name string) (Adapter, error)
func (r *Registry) List() []string
```

### Acceptance Criteria

- [ ] Claude Code adapter output là valid markdown
- [ ] Cursor adapter output có đúng `.mdc` frontmatter format
- [ ] Gemini CLI adapter output là valid markdown
- [ ] OpenCode adapter output có đúng frontmatter với `mode: subagent`
- [ ] Registry quản lý tất cả adapters
- [ ] Unit tests cho mỗi adapter

---

## Phase 5: Web UI (Embedded)

> **Goal**: Dashboard local để quản lý agents trực quan.

### 5.1 — REST API (`internal/web/api.go`)

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/config` | Lấy config hiện tại |
| `PUT` | `/api/config` | Cập nhật config |
| `GET` | `/api/agents` | Liệt kê tất cả agents |
| `GET` | `/api/agents/:name` | Lấy agent detail |
| `POST` | `/api/agents` | Tạo agent mới |
| `PUT` | `/api/agents/:name` | Cập nhật agent |
| `DELETE` | `/api/agents/:name` | Xóa agent |
| `GET` | `/api/targets` | Liệt kê targets |
| `GET` | `/api/preview/:agent/:target` | Preview rendered output |
| `POST` | `/api/sync` | Trigger sync |
| `POST` | `/api/sync/:target` | Trigger sync cho 1 target |
| `POST` | `/api/clean` | Trigger clean |
| `GET` | `/api/status` | Lấy status tổng quan |

### 5.2 — Web UI Design

**Layout**: Single-page dashboard với sidebar navigation.

**Pages:**

1. **Dashboard (Home)**
   - Tổng quan: số agents, số targets, trạng thái sync
   - Quick actions: Sync All, Clean
   - Recent activity

2. **Agents Manager**
   - Danh sách agents dạng cards
   - Click vào agent → mở editor (textarea cho markdown)
   - Tạo agent mới (modal/form)
   - Xóa agent (confirm dialog)
   - Tag hiển thị targets đã assign

3. **Targets Config**
   - Bảng matrix: agents × targets
   - Toggle checkbox để assign/unassign agent cho target
   - Visual indicator cho trạng thái sync

4. **Preview**
   - Chọn agent + target → xem rendered output
   - Side-by-side: source markdown vs rendered output
   - Syntax highlighting

5. **Status**
   - Trạng thái từng symlink (synced, out-of-date, broken)
   - Sync/Clean buttons per target

### 5.3 — UI Tech

- **HTML/CSS/JS** thuần — embedded via `go:embed`
- **CSS**: Dark mode, modern design, glassmorphism
- **JS**: Vanilla JS, fetch API cho REST calls
- **Markdown editor**: Simple textarea với preview (có thể dùng lightweight lib)
- **No build step** — files served trực tiếp

```go
//go:embed static/*
var staticFiles embed.FS
```

### 5.4 — HTTP Server (`internal/web/server.go`)

```go
type Server struct {
    engine   *engine.Engine
    config   *config.Config
    port     int
    rootDir  string
}

func NewServer(rootDir string, port int) *Server
func (s *Server) Start() error
```

### Acceptance Criteria

- [ ] `dotagen serve` mở dashboard tại `localhost:7890`
- [ ] CRUD agents hoạt động qua web UI
- [ ] Target mapping matrix hoạt động
- [ ] Preview render đúng format cho mỗi platform
- [ ] Sync/Clean trigger từ web UI
- [ ] Status hiển thị chính xác
- [ ] Dark mode, responsive design

---

## Phase 6: Polish & Release

> **Goal**: Production-ready, documented, installable.

### Tasks

#### 6.1 — Error Handling & Edge Cases

- `.dotagen/` không tồn tại → hướng dẫn chạy `dotagen init`
- Agent file rỗng → skip với warning
- Symlink target directory không tồn tại → auto-create
- Permission errors → clear error message
- Config YAML syntax error → hiển thị line number lỗi

#### 6.2 — README.md

- Logo/banner
- Quick start guide
- Detailed usage cho mỗi command
- Config reference
- Platform support matrix
- Screenshots web UI

#### 6.3 — Makefile

```makefile
build:        go build -o bin/dotagen ./cmd/dotagen
test:         go test ./...
lint:         golangci-lint run
install:      go install ./cmd/dotagen
dev:          go run ./cmd/dotagen serve
clean:        rm -rf bin/
```

#### 6.4 — Release

- Goreleaser config cho cross-platform builds
- GitHub Actions CI/CD
- `go install github.com/k0walski/dotagen@latest`

### Acceptance Criteria

- [ ] README hoàn chỉnh
- [ ] `make build` + `make test` pass
- [ ] Cross-platform build (Linux, macOS, Windows)
- [ ] `go install` hoạt động

---

## Appendix

### A. Platform Config Reference

#### Claude Code

```
Location: .claude/agents/{name}.md
Format:   Pure markdown
Docs:     https://docs.anthropic.com/en/docs/agents
```

#### Cursor

```
Location: .cursor/rules/{name}.mdc
Format:   YAML frontmatter + markdown
Fields:   description, globs, alwaysApply
Docs:     https://docs.cursor.com/context/rules
```

#### Gemini CLI

```
Location: .gemini/agents/{name}.md
Format:   Pure markdown
Docs:     https://github.com/google-gemini/gemini-cli
```

#### OpenCode

```
Location: .opencode/agents/{name}.md
Format:   YAML frontmatter + markdown
Fields:   description, mode (subagent)
Docs:     https://opencode.ai/docs/agents
```

### B. Example `.dotagen/` Setup

```
.dotagen/
├── config.yaml
├── agents/
│   ├── review-code.md
│   ├── planning.md
│   └── testing.md
├── .generated/          ← git-ignored
│   ├── claude-code/
│   │   ├── review-code.md
│   │   └── testing.md
│   ├── cursor/
│   │   ├── review-code.mdc
│   │   ├── planning.mdc
│   │   └── testing.mdc
│   ├── gemini-cli/
│   │   └── testing.md
│   └── opencode/
│       ├── planning.md
│       └── testing.md
└── .gitignore
```

**config.yaml:**

```yaml
targets:
  - claude-code
  - cursor
  - gemini-cli
  - opencode

agents:
  review-code:
    targets:
      - claude-code
      - cursor
  planning:
    targets:
      - cursor
      - opencode
  testing:
    targets: all
```

**agents/review-code.md:**

```markdown
---
description: Senior code reviewer
---

# Code Review Agent

You are a senior code reviewer with 10+ years of experience.

## Responsibilities

- Review code for security vulnerabilities
- Check for performance issues
- Ensure code readability and maintainability
- Verify test coverage

## Guidelines

- Be constructive, not critical
- Provide specific suggestions with code examples
- Prioritize issues by severity: 🔴 Critical > 🟡 Warning > 🔵 Info
```

### C. Implementation Priority

| Priority | Phase | Estimated Effort |
|---|---|---|
| 🔴 P0 | Phase 1: Scaffold | 1-2 hours |
| 🔴 P0 | Phase 2: Core Engine | 4-6 hours |
| 🔴 P0 | Phase 3: CLI Commands | 3-4 hours |
| 🔴 P0 | Phase 4: Platform Adapters | 2-3 hours |
| 🟡 P1 | Phase 5: Web UI | 6-8 hours |
| 🟢 P2 | Phase 6: Polish & Release | 2-3 hours |
| | **Total** | **~18-26 hours** |

> **Recommended approach**: Hoàn thành Phase 1→4 trước (CLI hoạt động end-to-end),
> sau đó Phase 5 (Web UI) có thể làm incremental.
