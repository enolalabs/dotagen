# Reviewer Mapping and Risk Selection

Use repository files and the reviewed plan/spec to infer the stack. Detection selects reviewer expertise; it does not justify adding more reviewers.

## Stack Signals

| Signal | Stack | Preferred technical reviewer |
|---|---|---|
| `go.mod` | Go | Go engineer |
| `package.json` | JavaScript/TypeScript | TypeScript/JavaScript engineer |
| `pyproject.toml`, `requirements.txt` | Python | Python engineer |
| `Cargo.toml` | Rust | Rust engineer |
| `*.csproj`, `*.sln` | .NET | C#/.NET engineer |
| `pom.xml`, `build.gradle` | Java/JVM | Java engineer |
| `Gemfile` | Ruby | Ruby engineer |
| `composer.json` | PHP | PHP engineer |
| `Package.swift` | Swift | Swift engineer |

Framework/database/cloud keywords refine the selected reviewer's context. They do not each receive a separate subagent.

## Group Triggers

| Group | Select when the reviewed scope contains |
|---|---|
| Technical | Always |
| Security | Auth, authorization, secrets, multi-tenancy, sensitive data, untrusted input, public endpoints |
| Performance | Explicit SLO/scale, concurrency, hot queries, large data, memory/resource bounds |
| Process | Deployment, migrations, rollback, lifecycle, CI/release gates, availability |
| Product | User-visible behavior, compatibility, rollout, or material scope/acceptance uncertainty |

For `Focused` mode, choose Technical plus the single highest-risk triggered group. If no specialist trigger is material, use Technical only. For `Differential`, use the owner of unresolved findings and add Technical only when contract consistency must be checked.

## Dispatch Rules

- Select one best-fit agent per group.
- Prefer a stack-specific technical agent over a generic architecture agent; do not dispatch both.
- Do not add database, cloud, frontend, or compliance agents merely because their technology is mentioned.
- Use a specialist only when the assigned scope contains a material question in that specialty.
- If a preferred agent is unavailable, use one general-purpose reviewer with the corresponding template.
- Respect the mode limits in `SKILL.md`.
