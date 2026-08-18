# 📚 Built-in Agents & Skills Catalog

> 🇻🇳 Vietnamese version: [CATALOG.vi.md](CATALOG.vi.md)

> **Scope:** This document describes the curated agent/skill bundles that ship with dotagen in detail: the **Game Studios** bundle (`dotagent-game-studios-*`, new) and the earlier `da-*`/`ds-*` catalog (VoltAgent + mattpocock). The full list of 842 built-in skills from 59 vendors, per-category counts, and how to enable `dotagent-*` entries are maintained in the [README](../README.md#built-in-skills).

> Every agent and skill is **disabled by default** — you choose what to enable, and for which platforms, in `config.yaml` or via the Web Dashboard.

---

## 🎮 Game Studios (49 agents · 73 skills)

The **Game Studios** bundle is imported from [donchitos/claude-code-game-studios](https://github.com/donchitos/claude-code-game-studios) (MIT) — a full "virtual game studio": engine specialists (Unity, Unreal, Godot), gameplay/engine/network/AI/UI programmers, designers, producer, QA, release, live-ops, narrative, audio and art roles, plus studio-style sprint / gate / release workflows.

- **Agents:** `dotagent-game-studios-<name>` (category `game-development`)
- **Skills:** directory `dotagent-game-studios-<name>/`, frontmatter `dotagent:game-studios:<name>` (category `Game Development`, vendor `game-studios`)
- Re-import from upstream with `python3 scripts/import-game-studios.py`.

> ⚠️ Many skills in this bundle assume the CCGS project layout (`design/`, `production/`, `.claude/rules/*`, templates under `.claude/docs/templates/*`). Run `/start` or `/adopt` first to bootstrap or audit that scaffolding.

### 🎯 Leadership & direction (5 agents)

| Agent | Description | When to use |
|---|---|---|
| **`creative-director`** | Highest creative authority: vision, tone, aesthetic direction, resolves design conflicts | Deciding the game vision, locking aesthetic direction, arbitrating design debates |
| **`technical-director`** | High-level technical decisions: engine architecture, technology choices, performance strategy, technical risk | Choosing engine/tech stack, overall architecture, technical risk assessment |
| **`lead-programmer`** | Code-level architecture, coding standards, code review, assigning work to specialist programmers | Code reviews, API design, refactoring strategy, translating design into code structure |
| **`producer`** | Sprint planning, milestone tracking, risk management, scope negotiation, cross-department coordination | Planning sprints, tracking milestones, handling scope creep |
| **`art-director`** | Visual identity: style guides, art bible, asset standards, palettes, art production pipeline | Art bible, asset standards, UI/UX visual review |

### 🎲 Game design (7 agents)

| Agent | Description | When to use |
|---|---|---|
| **`game-designer`** | Mechanical and systems design: core loops, progression, combat, economy, player experience | Designing/evaluating gameplay mechanics, writing GDDs |
| **`systems-designer`** | Detailed design of specific subsystems: combat formulas, progression curves, crafting, status effects | Mathematical/formula-level spec for one system |
| **`economy-designer`** | Resource economies, loot, progression curves, in-game markets | Loot tables, faucets/sinks, economy balance |
| **`level-designer`** | Spatial design, encounter layout, pacing, environmental storytelling | Level layouts, pacing plans, environmental storytelling guides |
| **`live-ops-designer`** | Post-launch content strategy: seasonal events, battle passes, cadence, retention | Live-service planning, events, retention mechanics |
| **`ux-designer`** | UX flows, interaction design, accessibility, information architecture, input handling | User flows, HUD, wireframes, input design |
| **`accessibility-specialist`** | Keeps the game playable by the widest audience; enforces accessibility standards, reviews UI | Accessibility audits, assist options (colorblind, remapping, subtitles…) |

### 💻 Programming (9 agents)

| Agent | Description | When to use |
|---|---|---|
| **`gameplay-programmer`** | Implements mechanics, player systems, combat, interactive features | Coding designed mechanics |
| **`engine-programmer`** | Core engine systems: rendering, physics, memory, resource loading, scene management | Engine/framework-level work |
| **`ai-programmer`** | Game AI: behavior trees, state machines, pathfinding, perception, NPC behavior | Enemy/NPC AI, decision-making systems |
| **`network-programmer`** | Multiplayer networking: state replication, lag compensation, matchmaking, protocols | Netcode, state sync, matchmaking |
| **`ui-programmer`** | UI systems: menus, HUD, inventory, dialogue boxes, UI framework | Implementing screens/HUD from UX specs |
| **`tools-programmer`** | Internal tools: editor extensions, content authoring tools, debug utilities, pipeline automation | Tools for designers/artists, pipeline automation |
| **`technical-artist`** | Bridges art and engineering: shaders, VFX, rendering optimization, art pipeline tools | Shaders/VFX, visual optimization, artist tooling |
| **`prototyper`** | Fast throwaway prototypes (concept prototype after brainstorm, feature prototype before GDD) | Proving an idea is fun before full design |
| **`devops-engineer`** | Build pipelines, CI/CD, version-control workflow, deployment | Build scripts, CI, branching workflow |

### 🎮 Engine specialists (15 agents)

| Agent | Description | When to use |
|---|---|---|
| **`unity-specialist`** | Authority on Unity patterns/APIs/optimization; MonoBehaviour vs DOTS | Any general Unity question |
| **`unity-dots-specialist`** | DOTS/ECS: Entity Component System, Jobs, Burst | Data-oriented performance in Unity |
| **`unity-shader-specialist`** | Shader Graph, HLSL, VFX Graph, URP/HDRP | Custom rendering in Unity |
| **`unity-ui-specialist`** | UI Toolkit (UXML/USS), UGUI, data binding, runtime UI performance | Building Unity UI |
| **`unity-addressables-specialist`** | Addressables: groups, load/unload, memory, catalogs, remote content | Asset/dynamic content management in Unity |
| **`unreal-specialist`** | Authority on Unreal patterns/APIs/optimization; Blueprint vs C++ | Any general Unreal question |
| **`ue-blueprint-specialist`** | Blueprint architecture, Blueprint/C++ boundary, graph optimization | Keeping Blueprints maintainable, deciding when to move to C++ |
| **`ue-gas-specialist`** | Gameplay Ability System: abilities, effects, attribute sets, tags, prediction | Combat/abilities built on GAS |
| **`ue-replication-specialist`** | Unreal networking: property replication, RPCs, prediction, relevancy, bandwidth | Multiplayer on Unreal |
| **`ue-umg-specialist`** | UMG/CommonUI: widgets, data binding, input routing, styling | Building Unreal UI |
| **`godot-specialist`** | Authority on Godot; GDScript vs C# vs GDExtension | Any general Godot question |
| **`godot-gdscript-specialist`** | GDScript: static typing, patterns, signals, coroutines, optimization | Writing/refactoring GDScript |
| **`godot-csharp-specialist`** | C# in Godot 4: .NET patterns, export attributes, signal delegates, async | Using C# in Godot |
| **`godot-gdextension-specialist`** | GDExtension: C/C++/Rust bindings, native optimization | Native code for Godot |
| **`godot-shader-specialist`** | Godot shading language, visual shaders, materials, particles, post-processing | Custom rendering in Godot |

### 📖 Narrative & content (4 agents)

| Agent | Description | When to use |
|---|---|---|
| **`narrative-director`** | Story architecture, world-building, characters, dialogue strategy | Story arcs, character development |
| **`world-builder`** | Detailed lore: factions, cultures, history, geography, ecology | Building the world and its rules |
| **`writer`** | Dialogue, lore entries, item descriptions, environmental text | Any player-facing text |
| **`localization-lead`** | i18n architecture, string tables, locale testing, translation pipeline | Multi-language support, translation workflow |

### 🔊 Audio (2 agents)

| Agent | Description | When to use |
|---|---|---|
| **`audio-director`** | Sonic identity: music direction, sound-design philosophy, implementation strategy, mix | Overall music/audio direction |
| **`sound-designer`** | Detailed SFX specs, audio events, mixing parameters | SFX spec sheets, audio event lists |

### 🧪 QA, performance, security, release & community (7 agents)

| Agent | Description | When to use |
|---|---|---|
| **`qa-lead`** | Test strategy, bug triage, quality gates, testing process | Test plans, severity assessment, release gates |
| **`qa-tester`** | Detailed test cases, bug reports, checklists | Test cases, regression checklists, bug reports |
| **`performance-analyst`** | Profiling, bottleneck identification, optimization recommendations, metric tracking | Slow game, performance budgets |
| **`security-engineer`** | Anti-cheat, exploits, data breaches; vulnerability review, save protection | Anti-cheat, save/network security |
| **`release-manager`** | Release pipeline: certification checklists, store submissions, platform requirements, versioning | Preparing a release, store submission |
| **`analytics-engineer`** | Telemetry, player-behavior tracking, A/B tests, data pipelines | Measuring/analyzing players |
| **`community-manager`** | Player-facing communication: patch notes, social, feedback collection, community bug triage | Writing patch notes, handling community feedback |

### 🛠 Skills — Getting started & discovery (6 skills)

| Skill | Description | When to use |
|---|---|---|
| **`start`** | First-time onboarding — asks where you are, then routes you to the right workflow | Brand new, unsure what to do |
| **`help`** | Analyzes what's done and your question, advises on the next step | "What should I do next?", feeling stuck |
| **`adopt`** | Brownfield onboarding: audits existing artifacts against templates, classifies gaps | Bringing an existing project into the CCGS workflow |
| **`onboard`** | Generates an onboarding doc for a new contributor/agent | Someone new joins the project |
| **`project-stage-detect`** | Detects project stage, gaps, and recommends next steps | "Where is the project right now?" |
| **`setup-engine`** | Configures engine + version, pins it in CLAUDE.md, populates engine reference docs | Right after brainstorm, before prototyping |

### 🎲 Skills — Concept & design (12 skills)

| Skill | Description | When to use |
|---|---|---|
| **`brainstorm`** | Guided ideation: from zero idea to a structured game concept document | Starting a new game |
| **`prototype`** | Concept prototype to validate the idea is worth designing before writing GDDs | Right after `/brainstorm` + `/setup-engine` |
| **`map-systems`** | Decomposes the concept into systems, maps dependencies, prioritizes design order, creates the systems index | After the concept, before GDDs |
| **`design-system`** | Guided, section-by-section GDD authoring for one system | Full GDD for a system |
| **`quick-design`** | Lightweight spec for small changes (tuning, minor mechanics, balance) | System GDD exists, only a tweak is needed |
| **`design-review`** | Reviews a GDD for completeness, consistency, implementability | Before handing a GDD to programming |
| **`review-all-gdds`** | Holistic cross-GDD review: contradictions, stale references | Before design milestones / pre-production |
| **`consistency-check`** | Scans GDDs against the entity registry for cross-document mismatches | After editing several GDDs |
| **`propagate-design-change`** | When a GDD changes, finds affected ADRs/traceability, produces a change report | Just revised an already-architected GDD |
| **`balance-check`** | Analyzes data files/formulas for outliers, broken progressions, degenerate strategies | After modifying balance data |
| **`ux-design`** | Section-by-section UX spec for a screen/flow/HUD | UX spec before UI work |
| **`ux-review`** | Validates a UX spec/HUD for completeness, accessibility, GDD alignment, implementation readiness | Before handing a UX spec to the UI programmer |

### 🎨 Skills — Art & assets (3 skills)

| Skill | Description | When to use |
|---|---|---|
| **`art-bible`** | Section-by-section Art Bible — the visual identity spec that gates all asset production | After brainstorm, before art production |
| **`asset-spec`** | Per-asset visual specs + AI generation prompts from GDDs/level docs/character profiles | Asset briefs for artists/AI |
| **`asset-audit`** | Audits assets for naming, size budgets, formats, pipeline; finds orphaned assets | Before milestones, after large asset imports |

### 🏛 Skills — Architecture & planning (7 skills)

| Skill | Description | When to use |
|---|---|---|
| **`create-architecture`** | Section-by-section master architecture document from all GDDs/ADRs/engine docs | After GDDs are approved |
| **`architecture-decision`** | Creates an ADR: context, alternatives, consequences | Any significant technical decision |
| **`architecture-review`** | Validates architecture against GDDs, builds a traceability matrix | After create-architecture, before create-epics |
| **`create-control-manifest`** | Flat rules sheet for programmers: must do / must never do, per system | After architecture is complete |
| **`create-epics`** | Translates GDDs + architecture into epics (one per module) | Start of production planning |
| **`create-stories`** | Breaks one epic into implementable story files with embedded GDD/ADR requirements | Before a sprint |
| **`reverse-document`** | Generates design/architecture docs from existing code | Under-documented projects |

### 🏃 Skills — Sprint & production (11 skills)

| Skill | Description | When to use |
|---|---|---|
| **`sprint-plan`** | Creates/updates a sprint plan from milestone, completed work, capacity | Sprint start |
| **`sprint-status`** | Fast sprint progress snapshot with burndown assessment | Mid-sprint |
| **`story-readiness`** | Validates a story is implementation-ready (GDD, ADR, engine notes, acceptance criteria) | Before starting a story |
| **`dev-story`** | Reads a story and implements it: loads full context, routes to the right programmer | Working a story |
| **`story-done`** | End-of-story review: checks each acceptance criterion, GDD/ADR deviations | Finishing a story |
| **`estimate`** | Effort estimate from complexity, dependencies, velocity, risk | Confidence-rated estimates |
| **`scope-check`** | Detects scope creep against the original plan | Suspected scope bloat |
| **`gate-check`** | Validates readiness to advance phases → PASS/CONCERNS/FAIL | Before moving to the next phase |
| **`milestone-review`** | Milestone review: completeness, quality, risk, go/no-go | End of milestone |
| **`retrospective`** | Sprint/milestone retro: velocity, blockers, patterns | End of sprint |
| **`vertical-slice`** | Production-quality end-to-end build to prove the full game loop | Pre-production |

### 🧑‍💻 Skills — Code & engineering (6 skills)

| Skill | Description | When to use |
|---|---|---|
| **`code-review`** | Architectural & quality review of a file or set of files | Before merging |
| **`tech-debt`** | Tracks, categorizes, prioritizes technical debt; maintains a debt register | Periodically / before milestones |
| **`perf-profile`** | Structured profiling workflow against budgets with prioritized recommendations | Missing performance targets |
| **`security-audit`** | Audits for save tampering, cheat vectors, network exploits, data exposure | Before release |
| **`localize`** | Full localization pipeline: hard-coded string scan, string tables, validation, translator briefs | Preparing multi-language support |
| **`content-audit`** | Compares GDD-specified content counts with implemented content | Tracking content progress |

### 🧪 Skills — Testing & QA (13 skills)

| Skill | Description | When to use |
|---|---|---|
| **`test-setup`** | Scaffolds the test framework + CI for the project's engine | Project start |
| **`test-helpers`** | Generates engine-specific test helper libraries | After test-setup |
| **`qa-plan`** | QA test plan for a sprint/feature | Sprint start or major feature |
| **`smoke-check`** | Critical-path smoke test gate before QA hand-off | Before QA hand-off |
| **`regression-suite`** | Maps coverage to GDD critical paths, finds fixed bugs lacking regression tests | Maintaining the regression suite |
| **`test-flakiness`** | Detects flaky tests from CI logs | Randomly red CI |
| **`test-evidence-review`** | Quality review of test files and manual evidence docs | Before gates/releases |
| **`soak-test`** | Soak-test protocol for extended play sessions | Hunting slow leaks/degradation |
| **`playtest-report`** | Playtest report template / analysis | After a playtest |
| **`bug-report`** | Structured bug report with full reproduction steps | Any bug |
| **`bug-triage`** | Re-evaluates all open bugs (priority vs severity), assigns to sprints | Periodic triage |
| **`skill-test`** / **`skill-improve`** | Validate & improve the CCGS skills themselves (linter, spec, audit) | Editing skills |

### 🚀 Skills — Release & live-ops (6 skills)

| Skill | Description | When to use |
|---|---|---|
| **`release-checklist`** | Pre-release checklist: build, certification, store metadata | Preparing a release |
| **`launch-checklist`** | Cross-department launch readiness + go/no-go | Before launch day |
| **`hotfix`** | Emergency fix workflow with audit trail, bypassing the sprint | Critical production bug |
| **`day-one-patch`** | Day-one patch: scope, prioritize, implement, QA gate | After gold, before launch |
| **`changelog`** | Internal + player-facing changelog from git/sprint/design data | End of sprint/release |
| **`patch-notes`** | Player-facing patch notes from git history/changelogs | Every update |

### 👥 Skills — Team orchestration (9 skills)

Each `team-*` skill orchestrates a group of agents through a complete pipeline.

| Skill | Agents | When to use |
|---|---|---|
| **`team-combat`** | game-designer, gameplay-programmer, ai-programmer, technical-artist, sound-designer, qa-tester | Design → implement → test a combat system |
| **`team-level`** | level-designer, narrative-director, world-builder, art-director, systems-designer, qa-tester | Building a complete area/level |
| **`team-narrative`** | narrative-director, writer, world-builder, level-designer | Story content & lore |
| **`team-ui`** | ux-designer, ui-programmer, art-director, qa-tester (integrates `/ux-design`, `/ux-review`) | UX spec → visual → implementation → polish |
| **`team-audio`** | audio-director, sound-designer, technical-artist, gameplay-programmer | Audio pipeline from direction to implementation |
| **`team-polish`** | performance-analyst, technical-artist, sound-designer, qa-tester | Optimize & polish a feature |
| **`team-qa`** | qa-lead, qa-tester | Full testing cycle |
| **`team-release`** | release-manager, qa-lead, devops-engineer, producer | Release candidate → deployment |
| **`team-live-ops`** | live-ops-designer, economy-designer, analytics-engineer, community-manager | Post-launch content planning |

---

## 🛠 Skills (Slash Commands) — mattpocock bundle

Skills are **structured workflows** (slash commands) an agent follows when triggered. Unlike an agent (which defines *who the agent is*), a skill defines *what the agent does and how*.

Each skill lives in a `dotagent-mattpocock-<name>/SKILL.md` directory (the `ds-*` names below are the legacy short names, now `dotagent-mattpocock-*`), optionally with a `references/` folder of supporting files.

### 🔧 Engineering (9 skills)

| Skill | Description | When to use |
|---|---|---|
| **`ds-diagnose`** | Diagnostic loop: Reproduce → Minimize → Hypothesize → Measure → Fix → Regression test | Hard bugs, perf regressions, "debug this" / "fix this" |
| **`ds-grill-with-docs`** | Challenges a plan against the domain model, sharpens terminology, updates CONTEXT.md and ADRs inline | Stress-testing a plan against existing docs and domain language |
| **`ds-improve-codebase-architecture`** | Finds architecture improvements based on CONTEXT.md and docs/adr/ | Refactoring, merging tight modules, making the codebase testable and AI-navigable |
| **`ds-setup-matt-pocock-skills`** | Sets up `## Agent skills` in AGENTS.md and `docs/agents/` for the repo | **Run first** — before `to-issues`, `triage`, `diagnose`, `tdd`, etc. |
| **`ds-tdd`** | Test-driven development with the red-green-refactor loop | Building a feature or fixing a bug via TDD, integration tests |
| **`ds-to-issues`** | Splits a plan/spec/PRD into independent issues as tracer-bullet vertical slices | Turning a plan into tracker tickets |
| **`ds-to-prd`** | Turns the current conversation context into a PRD | Producing a PRD from an ongoing conversation |
| **`ds-triage`** | Triages issues via a state machine with triage roles | Creating issues, classifying bugs/feature requests, preparing issues for autonomous agents |
| **`ds-zoom-out`** | Asks the agent for wider context and a high-level view | Unfamiliar code, need to see where it fits in the bigger picture |

### 🎯 Productivity (3 skills)

| Skill | Description | When to use |
|---|---|---|
| **`ds-caveman`** | Simplifies explanations to the bare minimum — "caveman" style | Quickly grasping a complex concept in plain language |
| **`ds-grill-me`** | Relentlessly interviews the user, resolving every branch of the decision tree | Stress-testing a plan/design, or when you say "grill me" |
| **`ds-write-a-skill`** | Creates a new skill with proper structure, progressive disclosure, bundled resources | Writing or building a new agent skill |

### 🔩 Misc (4 skills)

| Skill | Description | When to use |
|---|---|---|
| **`ds-git-guardrails`** | Sets up hooks blocking dangerous git commands (push, reset --hard, clean, branch -D) | Protecting the repo from destructive git in Claude Code |
| **`ds-migrate-to-shoehorn`** | Migrates test files from `as` type assertions to @total-typescript/shoehorn | TypeScript projects replacing `as` in tests with partial test data |
| **`ds-scaffold-exercises`** | Creates exercise directory structure: sections, problems, solutions, explainers | Scaffolding exercises for a course or new section |
| **`ds-setup-pre-commit`** | Sets up Husky pre-commit hooks with lint-staged, type checking, tests | Adding pre-commit hooks, Husky, or automatic checks on commit |

---

## 🤖 Agents — VoltAgent bundle

Agents are **virtual specialists** with a specific role and expertise. Each agent holds a system prompt guiding the AI to act in that role when invoked.

Each agent is a `dotagent-voltagent-<name>.md` file with YAML frontmatter (description, category) and Markdown content. Names in the tables are the short form.

### 💼 Business & Product (12 agents)

Agents supporting business operations, product management, and strategy.

| Agent | Description | When to use |
|---|---|---|
| **`business-analyst`** | Business process analysis, requirements gathering | Requirements analysis, process improvement, BRDs |
| **`content-marketer`** | Content strategy, SEO marketing, multi-channel campaigns | Content planning, SEO writing, marketing campaigns |
| **`customer-success-manager`** | Customer health scoring, retention strategy | Churn analysis, retention optimization, upsell opportunities |
| **`legal-advisor`** | Contract drafting, compliance review, IP protection | Contract review, legal risk assessment, compliance advice |
| **`license-engineer`** | OSI license selection, dependency compliance | Choosing a project license, checking dependency compliance |
| **`marketing-analyst`** | Campaign performance, attribution models | Campaign ROI analysis, growth modeling |
| **`product-manager`** | Feature prioritization, roadmaps, stakeholder coordination | Product planning, backlog prioritization, user stories |
| **`sales-engineer`** | Technical pre-sales, solution architecture, PoCs | Technical demos, customer solution design, PoCs |
| **`scrum-master`** | Sprint planning, retrospectives, impediment removal | Running Agile ceremonies, improving velocity, unblocking |
| **`technical-writer`** | API docs, user guides, SDK docs | Writing/improving technical docs, API docs, getting-started guides |
| **`ux-researcher`** | User research, usability testing, personas | UX surveys, user behavior analysis, design validation |
| **`wordpress-master`** | WordPress architecture, WooCommerce, security | Building/optimizing WordPress, custom themes/plugins, hardening |

### ⚙️ Core Development (11 agents)

Agents for core software development — backend, frontend, APIs, databases.

| Agent | Description | When to use |
|---|---|---|
| **`api-designer`** | API specs, RESTful patterns, GraphQL schemas | Designing new APIs, contract review, OpenAPI specs |
| **`backend-developer`** | Server-side APIs, microservices, robust backends | Building backends, service design, server-side optimization |
| **`database-architect`** | Schema design, query optimization, migrations | DB design, migrations, slow query tuning |
| **`frontend-developer`** | React/Vue/Angular frontends, responsive design | Building UIs, state management, UX optimization |
| **`fullstack-developer`** | End-to-end application development | Working across frontend and backend at once |
| **`graphql-developer`** | GraphQL schemas, resolvers, federation | Designing/implementing GraphQL APIs, federation setup |
| **`legacy-modernizer`** | Legacy system modernization | Migrating old systems, strangler fig, rewrites |
| **`low-level-designer`** | Class-level OOP/functional design, SOLID | Class diagrams, applying design patterns |
| **`microservices-architect`** | Distributed systems, service mesh, event-driven | Splitting monoliths, microservices design, event sourcing |
| **`ui-designer`** | Visual UI, design systems, component libraries | UI design, design systems, component libraries |
| **`websocket-engineer`** | Realtime bidirectional communication at scale | Realtime features: chat, notifications, live updates |

### 🧠 Data & AI (13 agents)

Agents for data science, machine learning, and AI.

| Agent | Description | When to use |
|---|---|---|
| **`ai-engineer`** | End-to-end AI systems, model selection, deployment | AI pipelines, model choice, model deployment |
| **`computer-vision-engineer`** | Image/video analysis, object detection, OCR | Image processing, object recognition, text from images |
| **`data-engineer`** | ETL pipelines, data warehouses, streaming | Data pipelines, warehouse design, streaming data |
| **`data-pipeline-architect`** | Large-scale data infrastructure, realtime processing | Data platform architecture, big data, realtime analytics |
| **`data-scientist`** | Statistical modeling, ML experiments, visualization | Data analysis, predictive models, A/B tests |
| **`data-visualization`** | Interactive dashboards, D3.js, Plotly | Charts, dashboards, visual reports |
| **`elasticsearch-specialist`** | Search clusters, query tuning, index management | Setting up/tuning Elasticsearch, full-text search |
| **`etl-specialist`** | Extract, transform, load pipelines | ETL jobs, data transformation, system sync |
| **`llm-architect`** | LLM applications, RAG, fine-tuning | AI apps/chatbots, RAG, model fine-tuning |
| **`ml-engineer`** | ML model development, training pipelines | Training models, ML pipelines, inference optimization |
| **`nlp-engineer`** | Text processing, sentiment analysis | Language analysis, text classification, entity extraction |
| **`playwright-expert`** | Browser automation, E2E testing, scraping | E2E tests, web scraping, browser automation |
| **`prompt-engineer`** | Prompt design, chain-of-thought, evaluation | Writing/optimizing prompts, LLM evaluation frameworks |

### 🔨 Developer Experience (14 agents)

Agents that improve developer experience and productivity.

| Agent | Description | When to use |
|---|---|---|
| **`build-engineer`** | Build performance, compilation optimization, scaling | Slow builds, CI build time optimization |
| **`cli-developer`** | Command-line tools and terminal apps | CLI tools, TUI apps, complex scripts |
| **`documentation-engineer`** | Documentation-as-code, architecture docs | Automated doc systems, docs sites, arch docs |
| **`git-specialist`** | Advanced git workflows, branching, history | Complex conflicts, rebases, git flow design |
| **`github-actions-specialist`** | CI/CD with GitHub Actions | Writing/optimizing GitHub Actions workflows |
| **`ide-plugin-developer`** | VS Code / JetBrains extensions | IDE plugins, code actions, language servers |
| **`json-wrangler`** | JSON/YAML transformation, validation, jq | Complex data transforms, JSON schemas, jq |
| **`monorepo-engineer`** | Nx/Turborepo/Lerna architecture | Monorepo setup, workspace management, task optimization |
| **`open-source-advisor`** | OSS contribution, community governance | OSS strategy, CONTRIBUTING.md, contributor management |
| **`refactoring-specialist`** | Code refactoring, tech-debt reduction | Poor-quality code, restructuring, reducing complexity |
| **`regex-master`** | Complex regex patterns, validation | Complex regex, text parsing, format validation |
| **`slack-expert`** | Slack apps, bots, API integration | Slack bots, webhooks, Slack apps |
| **`tooling-engineer`** | Developer tools, code generators | Dev tools, scaffolding, code gen |
| **`vibe-coder`** | Rapid prototyping, creative coding | Quick MVPs, prototypes, idea experiments |

### 🏗 Infrastructure (16 agents)

Agents for infrastructure, DevOps, and system administration.

| Agent | Description | When to use |
|---|---|---|
| **`azure-infra-engineer`** | Azure infrastructure, networking, deployment | Deploying/managing Azure infrastructure |
| **`cicd-engineer`** | CI/CD pipelines, deployment automation | Designing CI/CD pipelines, release automation |
| **`cloud-architect`** | Cloud infrastructure design, multi-cloud | Cloud architecture, cloud service selection |
| **`devops-engineer`** | Infrastructure automation, monitoring | IaC, monitoring stacks, deployment automation |
| **`docker-expert`** | Container optimization, multi-stage builds | Dockerfiles, image size, Compose |
| **`gcp-specialist`** | GCP services and architecture | Deploying on Google Cloud |
| **`kubernetes-specialist`** | K8s cluster management, Helm, operators | K8s deployments, Helm charts, custom operators |
| **`linux-sysadmin`** | Linux server administration, shell scripting | Linux servers, bash scripts |
| **`network-engineer`** | Network architecture, security, troubleshooting | Network design, firewalls, connectivity debugging |
| **`nginx-specialist`** | Nginx configuration, load balancing, reverse proxy | Nginx, SSL, reverse proxies |
| **`powershell-admin`** | Windows automation, Active Directory | PowerShell scripts, AD, GPO |
| **`security-engineer`** | Zero-trust architecture, CI/CD security | Security architecture, shift-left security |
| **`sre-engineer`** | SLO/SLI, error budgets, chaos engineering | Defining SLOs, reducing toil, chaos tests |
| **`terraform-engineer`** | Infrastructure as code, Terraform | Writing/refactoring Terraform, state management |
| **`terragrunt-expert`** | Terragrunt orchestration, DRY configuration | Terragrunt wrappers, multi-env deployment |
| **`windows-infra-admin`** | Windows Server, AD, Group Policy | Windows Server, AD, DNS, DHCP |

### 💬 Language Specialists (30 agents)

Specialists for individual programming languages and frameworks.

| Agent | Description | When to use |
|---|---|---|
| **`angular-architect`** | Angular 15+ enterprise | Large Angular apps, lazy loading, enterprise patterns |
| **`astro-developer`** | Astro framework, content-driven sites | Static/content-driven sites with Astro |
| **`cpp-systems-developer`** | C++ systems, memory management | C++ systems programming, performance, RAII |
| **`csharp-dotnet-developer`** | C#/.NET enterprise | .NET apps, ASP.NET Core, Blazor |
| **`django-developer`** | Django web and REST APIs | Django sites, DRF APIs, admin |
| **`elixir-phoenix-developer`** | Elixir/Phoenix realtime | Realtime apps, LiveView, distributed systems |
| **`flutter-developer`** | Flutter cross-platform | Mobile/web apps with Flutter, Dart |
| **`golang-pro`** | Go, concurrency, performance | Go apps, goroutines, performance tuning |
| **`java-enterprise-architect`** | Java enterprise, Spring | Large Java systems, Spring Boot, microservices |
| **`kotlin-expert`** | Kotlin/Android, coroutines | Android apps, Kotlin Multiplatform |
| **`laravel-expert`** | Laravel PHP, Eloquent ORM | Laravel sites, Livewire, PHP API backends |
| **`nestjs-architect`** | NestJS enterprise backend | NestJS APIs, TypeScript microservices |
| **`nextjs-developer`** | Next.js full-stack, SSR/SSG | Next.js sites, App Router, Server Components |
| **`nuxt-specialist`** | Nuxt 3, Vue ecosystem | Nuxt sites, auto-imports, server routes |
| **`perl-modernizer`** | Perl modernization, Moose | Modernizing legacy Perl codebases |
| **`php-engineer`** | PHP 8+, frameworks, performance | Modern PHP, Composer, performance tuning |
| **`python-pro`** | Python, async, data processing | Python apps, FastAPI, data processing |
| **`r-statistician`** | R statistics, data analysis | Statistical analysis, ggplot2, Shiny dashboards |
| **`rails-developer`** | Ruby on Rails | Rails sites, ActiveRecord, Hotwire |
| **`react-native-developer`** | React Native mobile | Cross-platform React Native apps |
| **`react-specialist`** | React 18+, hooks, state | React UIs, complex hooks, state management |
| **`ruby-pro`** | Ruby, metaprogramming | Plain Ruby, gems, DSLs |
| **`rust-engineer`** | Rust systems, memory safety | Rust, ownership, async, systems code |
| **`spring-boot-engineer`** | Spring Boot 3+ enterprise | Spring Boot apps, JPA, Security, Cloud |
| **`sql-pro`** | SQL, schemas, indexing | Complex query tuning, schema design, indexes |
| **`swift-expert`** | Swift/iOS/macOS native | iOS/macOS apps, SwiftUI, async/await |
| **`symfony-specialist`** | Symfony 6+/7+, Doctrine ORM | Symfony apps, Messenger, API Platform |
| **`typescript-pro`** | Advanced TypeScript types | Complex generics, type-level programming |
| **`vue-expert`** | Vue 3, Composition API | Vue apps, Composition API, Pinia |
| **`wordpress-master`** | Themes, plugins, WooCommerce | Custom WordPress, Gutenberg blocks |

### 🎭 Meta-Orchestration (11 agents)

Agents that coordinate and manage other agents.

| Agent | Description | When to use |
|---|---|---|
| **`agent-installer`** | Discover, browse, install agents | Finding and installing new agents for Claude Code |
| **`agent-organizer`** | Organizes and optimizes multi-agent teams | Assembling agent teams for complex projects |
| **`codebase-orchestrator`** | Repo-wide refactor governance | Large multi-file refactors with approval loops |
| **`context-manager`** | Shared state management between agents | Multiple agents need to sync data |
| **`error-coordinator`** | Coordinated error handling across components | Distributed errors across many components |
| **`it-ops-orchestrator`** | Multi-domain IT operations orchestration | Coordinating PS, AD, network domains at once |
| **`knowledge-synthesizer`** | Extracts patterns from agent interactions | Insights from agent activity history |
| **`multi-agent-coordinator`** | Coordinates concurrently running agents | Agents that must communicate and share state |
| **`performance-monitor`** | Observability infrastructure, metric tracking | System performance monitoring, anomaly detection |
| **`task-distributor`** | Task distribution, queue management, load balancing | Splitting work across agents/workers |
| **`workflow-orchestrator`** | Business workflow design | Multi-step automation with error handling |

### 🛡 Quality & Security (16 agents)

Agents for testing, security, and quality assurance.

| Agent | Description | When to use |
|---|---|---|
| **`accessibility-tester`** | WCAG compliance, accessibility testing | Checking/fixing accessibility, meeting WCAG |
| **`ad-security-reviewer`** | Active Directory security assessment | AD audits, privilege-escalation checks |
| **`ai-writing-auditor`** | Audits AI-generated content | Detecting/rewriting AI-generated text |
| **`architect-reviewer`** | System design review | Evaluating architectural decisions, design review |
| **`chaos-engineer`** | Controlled failure experiments | Resilience testing, failure injection, game days |
| **`code-reviewer`** | Comprehensive code review — security, quality | PR review, code quality or security checks |
| **`compliance-auditor`** | Regulatory compliance, audit controls | SOC2, GDPR, ISO |
| **`debugger`** | Bug diagnosis, root cause analysis | Hard bugs, root cause analysis |
| **`error-detective`** | Error correlation, failure-chain analysis | Cross-service errors, tracing failure chains |
| **`penetration-tester`** | Security penetration testing | Pentests, finding vulnerabilities before release |
| **`performance-engineer`** | Performance bottleneck identification | Slow apps, profiling, load tests |
| **`powershell-security-hardening`** | PowerShell security, hardening | Hardening PS remoting, constrained mode |
| **`qa-expert`** | QA strategy, test plans, coverage | Test strategy, coverage assessment |
| **`security-auditor`** | Security audits, compliance assessment | Comprehensive audits, security risk assessment |
| **`test-automator`** | Test automation frameworks, CI/CD | Test frameworks, testing in CI |
| **`ui-ux-tester`** | UI/UX functional testing | UI testing, user flows, defect finding |

### 🔬 Research & Analysis (8 agents)

Agents for research, analysis, and insight gathering.

| Agent | Description | When to use |
|---|---|---|
| **`competitive-analyst`** | Competitor analysis, benchmarking | Comparing with competitors, market positioning |
| **`data-researcher`** | Multi-source data discovery and validation | Collecting data from many sources, validating reliability |
| **`market-researcher`** | Market analysis, consumer behavior | Understanding markets, target customers, sizing |
| **`project-idea-validator`** | Idea stress-testing, competitor analysis | Checking feasibility of a new idea before investing |
| **`research-analyst`** | Multi-source research synthesis | Synthesizing information into reports |
| **`scientific-literature-researcher`** | Scientific literature, structured data | Finding/extracting data from scientific papers |
| **`search-specialist`** | Advanced search strategy, query optimization | Finding precise information quickly across sources |
| **`trend-analyst`** | Emerging trends, industry shift forecasting | Trend forecasting, future scenarios |

### 🌐 Specialized Domains (13 agents)

Agents for specialized fields.

| Agent | Description | When to use |
|---|---|---|
| **`api-documenter`** | API documentation, OpenAPI specs | Writing/improving API docs, Swagger specs |
| **`blockchain-developer`** | Smart contracts, DApps, blockchain | Smart contracts, DApps, DeFi protocols |
| **`embedded-systems`** | Firmware, RTOS, microcontrollers | Embedded programming, firmware, RTOS |
| **`fintech-engineer`** | Payment systems, financial compliance | Payment systems, PSD2, banking APIs |
| **`game-developer`** | Game systems, graphics, multiplayer | Games, rendering, physics, netcode |
| **`healthcare-admin`** | Healthcare administration, HIPAA compliance | Healthcare projects, EMR, HIPAA |
| **`iot-engineer`** | IoT device management, edge computing | IoT platforms, device management, edge |
| **`m365-admin`** | Microsoft 365 administration automation | Automating Exchange, SharePoint, Teams |
| **`mobile-app-developer`** | iOS/Android apps | Native or cross-platform mobile apps |
| **`payment-integration`** | Payment gateway integration, PCI | Stripe, PayPal, PCI compliance |
| **`quant-analyst`** | Quantitative trading, financial modeling | Trading strategies, risk models, backtesting |
| **`risk-manager`** | Risk identification/mitigation | Risk assessment, control frameworks |
| **`seo-specialist`** | SEO audits, keyword strategy | SEO optimization, technical audits, keyword research |

---

## 🚀 Usage

### Enable agents/skills in config

```yaml
# .dotagen/config.yaml
agents:
  dotagent-game-studios-lead-programmer:
    targets: all          # Enable for every platform
  dotagent-voltagent-code-reviewer:
    targets:
      - claude-code       # Claude Code only

skills:
  dotagent-game-studios-gate-check:
    targets: all
  dotagent-superpowers-test-driven-development:
    targets:
      - claude-code
      - cursor
```

### Sync

```bash
dotagen sync              # Sync everything
dotagen sync claude-code  # Sync Claude Code only
```

### Manage via CLI

```bash
dotagen skill list                 # List skills
dotagen skill create my-workflow   # Create a new skill
dotagen status                     # Check status
```

### Manage via Web Dashboard

```bash
dotagen serve   # Open the dashboard at http://localhost:7890
```

---

## 📝 Sources

- **49 agents + 73 skills** Game Studios from [donchitos/claude-code-game-studios](https://github.com/donchitos/claude-code-game-studios) (MIT)
- **154 agents** `dotagent-voltagent-*` from [VoltAgent/awesome-claude-code-subagents](https://github.com/VoltAgent/awesome-claude-code-subagents)
- **26 skills** `dotagent-mattpocock-*` from [mattpocock/skills](https://github.com/mattpocock/skills)
- **14 skills** `dotagent-superpowers-*` from [obra/superpowers](https://github.com/obra/superpowers)
- The remaining vendor skills come from the [awesome-agent-skills](https://github.com/enolalabs/awesome-agent-skills) registry; [Hallmark](https://github.com/Nutlope/hallmark) is a separately pinned MIT snapshot
