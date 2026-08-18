---
name: "dotagent:enolalabs:spec-review"
description: "Run a risk-based, token-efficient review of a product or architecture spec. Defaults to focused review and uses differential re-review after changes."
category: "Developer Tools"
vendor: "enolalabs"
---

# Spec Review

Review a spec for decision completeness and material delivery risk. Prefer the smallest review scope that can answer: **is the problem, boundary, and acceptance contract clear enough to plan?**

## When to Use

Use after brainstorming produces a spec and before writing an implementation plan.

```text
brainstorming → spec-review → writing-plans → plan-review → implementation
```

## Review Modes

| Mode | Use when | Default reviewers | Context |
|---|---|---|---|
| **Focused** | First review of a normal spec or phase | Technical + at most one highest-risk specialist | Target phase/sections and directly relevant repository context |
| **Full** | Explicit request, or a one-time release/foundation contract whose cross-cutting risks cannot be covered phase-by-phase | All relevant groups, at most one reviewer per group | Entire spec and required supporting docs |
| **Differential** | Any re-review after edits | Technical + only specialists that own unresolved findings, maximum two reviewers | Changed sections plus unresolved findings |

`Focused` is the default for an initial review. `Differential` is the default whenever the user says **re-review**, **review again**, or asks whether a revised spec is ready. Never repeat stack or reviewer-selection questions during the same review chain unless the scope materially changes.

## Process

### 1. Locate and Bound the Review

- Use an explicit path from the user when present.
- Otherwise scan `docs/superpowers/specs/`. Ask only when multiple plausible targets remain.
- If the spec contains phases, default to the phase being planned next plus global invariants it depends on.
- Read only supporting documents directly referenced by the reviewed sections.
- Record the target path, reviewed sections, mode, and current commit.

For a differential review, locate the prior report and comparison base. Prefer the report's recorded base, then the pre-fix commit, current worktree diff, or sections identified in conversation. If no reliable diff exists, inspect only unresolved findings and their referenced sections. Do not silently fall back to a full review.

### 2. Detect Context and Select Risk Reviewers

Follow [subagent-mapping.md](subagent-mapping.md). Infer the stack and product type from the repository and spec. Ask for confirmation only if uncertainty would change a material conclusion. Reuse prior answers.

Technical review is the default anchor. Select at most one additional group for the highest current risk:

- Product: user-visible flows, business outcomes, compatibility, or scope uncertainty.
- Security: authentication, authorization, secrets, multi-tenancy, untrusted input, or sensitive data.
- Performance: explicit scale/SLOs, concurrency, large data, or resource bounds.
- Process: deployment, migrations, rollback, lifecycle, availability, or release gates.

Respect an explicit user choice of mode or groups. Otherwise do not ask a separate reviewer-selection question.

### 3. Apply the Spec Sufficiency Standard

A spec is ready for planning when it makes material decisions explicit enough that plan authors do not have to invent product or architecture policy. As relevant to its scope, it should define:

- problem, users/actors, goals, and non-goals;
- user or system flows and observable outcomes;
- boundaries, ownership, dependencies, and durable invariants;
- failure behavior and security/trust boundaries;
- performance/operational constraints where they affect architecture;
- measurable acceptance criteria;
- true open decisions, with owner or decision point.

A spec should not be required to contain implementation file paths, full type definitions, SQL, production code, exhaustive test bodies, or library-level mechanics unless those details are themselves architectural contracts. Push implementation choices into the plan rather than expanding the spec.

An ambiguity is a finding only when two reasonable interpretations would change behavior, compatibility, architecture, risk, or acceptance. Missing generic sections are not automatically Critical.

### 4. Dispatch the Minimum Useful Reviewers

Read the selected `reviewers/<group>.md` templates and dispatch only when subagents are available and allowed.

- Focused: maximum two reviewers.
- Differential: maximum two reviewers; one is preferred for localized edits.
- Full: maximum one reviewer per selected group.
- Use one best-fit agent per group, not a generic and specialist reviewer for the same concern.

When reviewers share the filesystem, pass the spec path, assigned sections, mode, comparison base, and unresolved finding IDs. Do not paste the full spec into every prompt. Reviewers should read only their scope and direct references.

If subagents are unavailable, perform the same bounded review directly.

### 5. Consolidate and Calibrate Findings

Every finding needs evidence, a precise section reference, impact, and the smallest adequate change. Deduplicate shared root causes before counting them.

Severity means:

- **Critical:** the core problem/actor/outcome is contradictory or infeasible; a required behavior is absent; trust/data boundaries permit a credible severe failure; or planning would require inventing a fundamental product/architecture decision.
- **Important:** a material ambiguity, missing acceptance criterion, likely rework, or meaningful security/performance/operational risk in the current scope.
- **Minor:** non-blocking clarity, polish, or optional hardening.

Do not report generic best practices, future-phase features, or missing details that properly belong in implementation planning.

Each reviewer may return at most three Critical/Important findings and two Minor findings. Ask for the highest-impact independent findings and omit low-confidence speculation. The consolidated report should normally contain no more than five blockers. If more independent blockers exist, state that the review is truncated, return `Needs Changes`, and review the remainder only after the first batch is fixed.

During differential review:

- classify prior blockers as `resolved`, `partially resolved`, or `open`;
- add a new blocker only when introduced/exposed by the edit or required to avoid an unsafe approval;
- do not reopen accepted product/architecture decisions without new evidence;
- do not repeat unchanged Minor findings.

### 6. Decide the Verdict

| Verdict | Condition |
|---|---|
| **Approved** | No blocking findings remain; the spec is ready for planning |
| **Approved with Follow-ups** | No blockers; only non-blocking follow-ups remain |
| **Needs Changes** | Localized blocking findings remain |
| **Rejected** | The problem, scope, or architecture needs fundamental rework |

Mark every Critical/Important finding `Blocking: yes|no`. Critical is normally blocking. Important blocks only when the issue must be decided before planning the reviewed scope.

### 7. Save a Compact Report

Save to:

```text
docs/superpowers/reviews/YYYY-MM-DD-<topic>-spec-review.md
```

Use this structure:

```markdown
# Spec Review: [Topic]

**Mode:** Focused | Full | Differential
**Scope:** [phase/sections]
**Spec:** `path`
**Base:** [commit/report, for differential review]
**Reviewers:** [groups/agents]

## Verdict: [verdict]

[One-paragraph summary]

## Blocking Findings

1. **[ID] [title]** — Critical|Important — Blocking: yes
   - Evidence: [section and exact conflict or omission]
   - Impact: [planning/delivery consequence]
   - Minimal fix: [bounded decision or text change]

## Follow-ups

[Only non-blocking items worth preserving]

## Coverage and Resolution

- [acceptance criterion or prior finding] → [section] → covered|resolved|open
```

Do not include empty reviewer sections, copied checklists, or long generic recommendations. Mention two or three concrete strengths at most.

### 8. Fix and Re-review Without Loops

Review is read-only by default. Do not edit or commit the spec because findings exist.

If the user asks to fix findings:

- batch clear fixes into one editing pass;
- ask only for decisions that affect scope, behavior, compatibility, architecture, or risk;
- group related choices into one concise question where possible;
- do not commit unless the user explicitly asks.

After fixes, run one `Differential` verification when requested or when the user asked to “fix and verify.” A second full review requires an explicit request or a material change to the problem, scope, trust boundary, or architecture. If blockers are gone, stop and declare the spec ready for planning.

## Key Principles

- Optimize for decision quality per token, not reviewer count.
- Review the phase being planned, not the whole roadmap repeatedly.
- Re-review edits and unresolved blockers, not unchanged text.
- Keep product/architecture decisions in the spec and implementation mechanics in the plan.
- Findings must be evidence-backed, deduplicated, and scope-relevant.
- Ask fewer questions and reuse prior decisions.
- Never auto-edit or auto-commit during a review-only request.
