---
name: "dotagent:enolalabs:plan-review"
description: "Run a risk-based, token-efficient review of an implementation plan. Defaults to focused review and uses differential re-review after changes."
category: "Developer Tools"
vendor: "enolalabs"
---

# Plan Review

Review an implementation plan for executability and material delivery risk. Prefer the smallest review scope that can answer: **can implementation begin safely, and what truly blocks it?**

## When to Use

Use after a plan is written and before executing the affected package or milestone.

```text
brainstorming → spec-review → writing-plans → plan-review → implementation
```

For a plan set, review global architecture and contracts once, then review each package immediately before execution. Do not repeatedly review the entire plan set.

## Review Modes

| Mode | Use when | Default reviewers | Context |
|---|---|---|---|
| **Focused** | First review of a normal plan or package | Technical + at most one risk specialist | Target plan, relevant spec sections, global constraints, direct dependencies |
| **Full** | Explicit user request, or a one-time release/foundation gate whose cross-cutting risks cannot be covered package-by-package | All relevant groups, at most one reviewer per group | Full plan set and spec |
| **Differential** | Any re-review after findings were fixed | Technical + only specialists that own unresolved findings, maximum two reviewers | Changed sections plus unresolved findings |

`Focused` is the default for an initial review. `Differential` is the default whenever the user says **re-review**, **review again**, or asks whether a revised plan is ready. Never ask the user to choose review groups again during the same review chain unless scope or risk changed materially.

## Process

### 1. Locate and Bound the Review

- Use an explicit path from the user when present.
- Otherwise scan `docs/superpowers/plans/`. Ask only when multiple plausible targets remain.
- If the target is an index referencing package plans, default to the selected package plus global constraints and direct dependency contracts. Review the whole set only in `Full` mode.
- Read the original spec only where needed for traceability. Do not load unrelated future-phase sections.
- Record the target path, scope, mode, and current commit.

For a differential review, also locate the previous review report and determine the comparison base. Prefer, in order:

1. base commit recorded in the previous report;
2. the commit immediately before the review fixes;
3. the current worktree diff;
4. changed sections identified from the conversation.

If no reliable diff exists, review only unresolved findings and the sections they reference. Do not silently fall back to a full review.

### 2. Detect Context Without Reconfirming the Obvious

Follow [subagent-mapping.md](subagent-mapping.md). Infer the stack from the plan, spec, and repository. Ask the user only when an ambiguity would change the review or recommended architecture. Reuse stack decisions already confirmed in the conversation or previous report.

Choose the specialist for the highest material risk:

- Security: authentication, authorization, secrets, multi-tenancy, untrusted input, or sensitive data.
- Performance: explicit SLOs, concurrency, large data, hot database paths, or resource limits.
- Process: migrations, deployment, rollback, lifecycle, or release gates.
- Product: user-visible flows, scope coverage, compatibility, or acceptance behavior.

Respect explicitly requested groups or mode. Otherwise do not ask a separate group-selection question.

### 3. Apply the Plan Sufficiency Standard

A plan is executable when it provides enough information to implement and verify the intended behavior without reopening architectural decisions. A task normally needs:

- concrete files or owned components;
- intended outcome and externally consumed contracts/invariants;
- dependency ordering where relevant;
- named test scenarios or assertions, including important failure cases;
- exact verification commands and package/milestone gates.

Representative signatures, SQL shapes, pseudocode, or snippets are useful when they freeze a boundary. Full production bodies and copy-pasteable test bodies are **not required**. TDD belongs in execution; the plan must specify behavior and evidence, not pre-implement the feature.

Treat a placeholder as blocking only when it leaves a required behavior, contract, security property, data invariant, or architectural choice undecided. Implementation-local details bounded by tests and contracts may be decided during execution.

Commit commands are optional unless immutable history or a named workflow gate depends on them.

### 4. Dispatch the Minimum Useful Reviewers

Read the selected `reviewers/<group>.md` templates and dispatch reviewers only if subagents are available and allowed.

- Focused: maximum two reviewers.
- Differential: maximum two reviewers; only one is preferred for a localized change.
- Full: maximum one reviewer per selected group.
- Use one best-fit agent per group; do not dispatch both a generic and a stack-specific agent for the same concern.

When reviewers share the filesystem, pass paths, scope, mode, comparison base, and unresolved finding IDs. Do not paste the full plan/spec into every prompt. Tell reviewers to read only the assigned scope and direct references.

If subagents are unavailable, perform the same bounded review directly; do not broaden scope to compensate.

### 5. Consolidate and Calibrate Findings

Every finding must include evidence, a precise plan reference, impact, and the smallest adequate fix. Deduplicate findings that share one root cause before counting them.

Severity means:

- **Critical:** implementation is impossible or contradictory; required spec behavior is absent; a cross-task contract cannot compile or compose; or the plan creates a credible security, isolation, or data-loss failure.
- **Important:** likely correctness failure, major rework, unverifiable acceptance criterion, or material operational/performance risk within the stated phase.
- **Minor:** non-blocking clarity, polish, or optimization.

Do not promote an issue because a code sample is incomplete. Do not report future-phase hardening, generic best practices, or preferences not required by the spec and current phase.

Each reviewer may return at most three Critical/Important findings and two Minor findings. Ask for the highest-impact independent findings; omit duplicates and low-confidence speculation. The consolidated report should normally contain no more than five blocking findings. If more independent blockers exist, state that the review is truncated, return `Needs Changes`, and review the remainder only after the first batch is fixed.

During differential review:

- verify each prior blocking finding as `resolved`, `partially resolved`, or `open`;
- report a new blocker only if it is introduced/exposed by the change or is a clear correctness/security issue that would make approval unsafe;
- do not reopen accepted decisions without new evidence;
- do not re-report unchanged Minor findings.

### 6. Decide the Verdict

| Verdict | Condition |
|---|---|
| **Approved** | No blocking findings remain |
| **Approved with Follow-ups** | No blockers; only explicitly non-blocking follow-ups remain |
| **Needs Changes** | Localized blocking findings remain |
| **Rejected** | A fundamental rewrite or missing architectural decision prevents a safe plan |

Mark each Critical/Important finding `Blocking: yes|no`. Critical findings are normally blocking. An Important finding is blocking only when it must be resolved before starting the reviewed scope.

### 7. Save a Compact Report

Save to:

```text
docs/superpowers/reviews/YYYY-MM-DD-<topic>-plan-review.md
```

Use this structure:

```markdown
# Plan Review: [Topic]

**Mode:** Focused | Full | Differential
**Scope:** [plan/package/sections]
**Plan:** `path`
**Spec:** `path` or selected sections
**Base:** [commit/report, for differential review]
**Reviewers:** [groups/agents]

## Verdict: [verdict]

[One-paragraph summary]

## Blocking Findings

1. **[ID] [title]** — Critical|Important — Blocking: yes
   - Evidence: [task/step and exact conflict or omission]
   - Impact: [implementation consequence]
   - Minimal fix: [bounded change]

## Follow-ups

[Only non-blocking items worth preserving]

## Coverage and Resolution

- [spec criterion or prior finding] → [task/section] → covered|resolved|open
```

Do not pad the report with empty group sections, repeated checklists, or long strengths lists. Mention two or three concrete strengths at most.

### 8. Fix and Re-review Without Loops

Review is read-only by default. Do not modify the plan or commit merely because findings exist.

If the user asks to fix findings:

- batch clear fixes into one editing pass;
- ask only for choices that change architecture, scope, compatibility, or risk;
- group related choices into one concise question where possible;
- do not commit unless the user explicitly asks for a commit.

After fixes, run one `Differential` verification when requested or when the user asked to “fix and verify.” A second full review requires an explicit user request or a material change to global architecture/contracts. If no blockers remain, stop and declare the reviewed scope ready for execution.

## Key Principles

- Optimize for decision quality per token, not reviewer count.
- Review the current execution boundary, not every future package.
- Re-review changes, not unchanged documents.
- A plan freezes behavior and boundaries; it does not contain the implementation.
- Findings must be evidence-backed, deduplicated, and phase-relevant.
- Ask fewer questions and reuse decisions already made.
- Never auto-edit or auto-commit during a review-only request.
