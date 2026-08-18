# Reviewer Mapping and Risk Selection

Use repository files and the reviewed spec to infer the stack and product type. Detection selects reviewer expertise; it does not justify adding more reviewers.

## Stack Signals

| Signal | Stack | Preferred technical reviewer |
|---|---|---|
| `go.mod` | Go | Go architect/engineer |
| `package.json` | JavaScript/TypeScript | TypeScript/JavaScript architect |
| `pyproject.toml`, `requirements.txt` | Python | Python architect/engineer |
| `Cargo.toml` | Rust | Rust engineer |
| `*.csproj`, `*.sln` | .NET | C#/.NET architect |
| `pom.xml`, `build.gradle` | Java/JVM | Java architect |
| `Gemfile` | Ruby | Ruby architect |
| `composer.json` | PHP | PHP architect |
| `Package.swift` | Swift | Swift engineer |

Framework, database, and cloud keywords refine the selected reviewer's context. They do not each receive a separate subagent.

## Group Triggers

| Group | Select when the reviewed scope contains |
|---|---|
| Technical | Always |
| Product | User-visible outcomes, business scope, compatibility, or acceptance uncertainty |
| Security | Auth, authorization, secrets, multi-tenancy, sensitive data, untrusted input, public endpoints |
| Performance | Explicit SLO/scale, concurrency, large data, or resource constraints that affect architecture |
| Process | Deployment, migrations, rollback, lifecycle, availability, or release gates |

For `Focused` mode, choose Technical plus the single highest-risk triggered group. For an internal architecture spec with no material product ambiguity, do not add Product automatically. For `Differential`, use the owner of unresolved findings and add Technical only when architecture consistency must be checked.

## Dispatch Rules

- Select one best-fit agent per group.
- Prefer a stack-specific technical agent over a generic architecture agent; do not dispatch both.
- Do not add database, cloud, frontend, or compliance agents merely because their technology is mentioned.
- Use a specialist only when the assigned scope contains a material question in that specialty.
- If a preferred agent is unavailable, use one general-purpose reviewer with the corresponding template.
- Respect the mode limits in `SKILL.md`.
