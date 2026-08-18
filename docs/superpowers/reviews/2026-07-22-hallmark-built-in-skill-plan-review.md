# Plan Review: Hallmark Built-in Skill

**Date:** 2026-07-22  
**Plan:** docs/superpowers/plans/2026-07-22-hallmark-built-in-skill.md  
**Spec:** docs/superpowers/specs/2026-07-22-hallmark-built-in-skill-design.md  
**Tech Stack:** Go 1.26.2 CLI/library, go:embed, frontmatter, SHA-256  
**Reviewers:** technical, performance, security, process & operations, product

---

## Verdict: Needs Changes

**Summary:** The plan has a sound red/green import sequence and carefully bounded artifact, but its core supply-chain ledger and automation requirements are not yet executable or complete enough to be a reliable guard for future snapshots.

**Finding counts:** 2 Critical, 6 Important, 4 Minor

**Spec coverage:** 10 of 10 requirements have tasks; 5 are partial because the plan needs stronger exact-set, CI, and documentation steps.

---

## Findings

### Technical quality

#### Critical (Must Fix)

1. **Metadata test is described but not supplied as executable code**
   - Task/Section: Task 1, Step 1
   - Issue: The plan says to add TestHallmarkSkillMetadata but provides neither a complete function nor package/import set.
   - Impact: A task-by-task implementer can omit or inconsistently implement a required metadata gate.
   - Recommended Fix: Include the complete copy-pasteable function, including frontmatter parsing and error cases.

2. **Manifest generation is not deterministic or executable**
   - Task/Section: Task 2, Step 4
   - Issue: The plan gives a schema and digest listing command but no exact local generator that emits valid sorted YAML for all 106 entries.
   - Impact: The central integrity artifact can be malformed, incomplete, or non-reproducible.
   - Recommended Fix: Specify an exact temporary local generator command/script and exact output schema.

#### Important (Should Fix)

1. **Manifest-declared file count and bytes are not asserted**
   - Task/Section: Task 1, Step 1
   - Recommended Fix: Compare decoded upstream_files/upstream_bytes to the manifest list/summed bytes before asserting fixed reviewed values.

2. **Manifest and embedded file tree are not compared for exact equality**
   - Task/Section: Task 1, Step 1; Task 1, Step 2
   - Impact: Unmanifested embedded content can evade the audit while still being copied by init.
   - Recommended Fix: Compare sorted manifest paths to ListSkillFiles after excluding SOURCE.md; validate every path is SKILL.md or a relative references/*.md path. Make the init test assert the same expected destination set.

3. **Absolute-link and malformed-manifest handling are incomplete**
   - Task/Section: Task 1, Step 1; Task 2, Step 2
   - Recommended Fix: use checked type assertions with index/key diagnostics and add a concrete scan/review criterion for absolute filesystem links and references.

### Security

#### Important (Should Fix)

1. **Safety review relies too heavily on a narrow regex**
   - Task/Section: Task 2, Step 2; SOURCE.md checklist
   - Impact: Unsafe instructions can evade literal terms yet be distributed to enabled model clients.
   - Recommended Fix: Require human sign-off after reviewing the complete pinned file/diff set against explicit remote-loading, secret/exfiltration, and unsafe-execution criteria; retain grep only as triage.

2. **Non-regular files and symlinks are not rejected**
   - Task/Section: Task 2, Steps 2–3
   - Impact: A future upstream symlink can escape the reviewed source boundary or break embedding.
   - Recommended Fix: reject symlinks and every entry that is not a regular Markdown file before copying; after copy, enforce the exact manifest set.

### Process & operations

#### Important (Should Fix)

1. **CI enforcement required by the spec has no task**
   - Task/Section: Task 4, Steps 1–2
   - Impact: Future updates can bypass the integrity and init tests.
   - Recommended Fix: add a task for the repository's authoritative CI pipeline to run go test ./... on pull requests and pushes, and verify it exists.

### Product

#### Important (Should Fix)

1. **README attribution remains inaccurate**
   - Task/Section: Task 3, Steps 1–2
   - Impact: README would claim the whole library comes from awesome-agent-skills and obra/superpowers despite the standalone Nutlope snapshot.
   - Recommended Fix: update Built-in Skills and Acknowledgments text to distinguish Hallmark's pinned MIT source and link to Nutlope/hallmark.

#### Minor (Nice to Have)

1. Record focused-test duration or binary-size delta in the final handoff.
2. Use a newly-created temporary import directory and remove it after verification.
3. Phrase Task 4 change-boundary checks to identify task-owned changes, not require a clean user worktree.
4. Add an optional manual enabled-path check: enable the initialized skill for codex, run the normal sync/list workflow, and confirm it renders.

---

## Spec Coverage Gaps

None. The gaps are implementation-quality gaps in requirements already mapped to Tasks 1–4.

## Strengths

- The plan follows an appropriate red-test, pinned import, documentation, and full-suite sequence.
- It preserves the existing embed/init architecture instead of adding unnecessary production code.
- Fixed source revision, 106-file/675085-byte bound, digest verification, and disabled default config are proportionate controls for this CLI scope.
- No load testing, database, deployment, or web security tasks are needed for a static local skill snapshot.

## Recommendations

Strengthen the manifest as an exact audited set, make its generation reproducible, and make CI the durable future-update gate. Update user-facing attribution alongside counts so the newly advertised source is honest.

## Resolution

The requester selected an in-repository GitHub Actions gate. The plan now
includes a pull-request/push workflow running `go test ./...`, full
copy-pasteable metadata and exact-set test requirements, deterministic ledger
generation, rejection of symlinks/non-Markdown files, complete-tree reviewer
sign-off, and corrected README attribution. It also adds an optional enabled
path acceptance check and preserves unrelated working-tree changes during
handoff.
