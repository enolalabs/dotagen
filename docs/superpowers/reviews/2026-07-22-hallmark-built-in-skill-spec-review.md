# Spec Review: Hallmark built-in skill

**Date:** 2026-07-22  
**Spec:** `docs/superpowers/specs/2026-07-22-hallmark-built-in-skill-design.md`  
**Tech Stack:** Go 1.26.2 CLI/library, `go:embed`, Markdown documentation  
**Reviewers:** technical architecture, performance, security, process & operations, product

---

## Verdict: Needs Changes

**Summary:** The offline, disabled-by-default embedding approach fits Dotagen's existing initialization model and is deliberately small in scope. The spec needs a reproducible, reviewable upstream snapshot contract plus precise metadata and size/test acceptance criteria before implementation planning.

**Finding counts:** 0 Critical, 4 Important, 5 Minor

---

## Findings

### Technical quality

#### Important (Should Fix)

1. **The upstream snapshot is not reproducibly identified**
   - Section: Snapshot rules
   - Issue: Importing from mutable `main` does not record the repository URL, source path, immutable commit SHA, or import date.
   - Impact: Maintainers cannot reproduce, audit, or reliably update the embedded prompt content.
   - Recommended fix: Require a committed provenance record for the pinned commit and imported path set, with a manual update and verification checklist.

2. **The reference-file test contract is underspecified**
   - Section: Validation
   - Issue: “Expected reference files” does not define the complete required path set.
   - Impact: A required linked reference could be omitted while a superficial test still passes.
   - Recommended fix: Test an explicit expected list derived from the pinned snapshot after both embedding and `runInit`.

#### Minor (Nice to Have)

1. **Incompatible-file handling is not actionable**
   - Section: Snapshot rules
   - Recommendation: State that import/PR validation fails with the offending path and reason when it finds an unsupported file.

2. **Documentation counts lack a source of truth**
   - Section: Documentation
   - Recommendation: Derive published counts from the final embedded skill list at implementation time.

### Performance

#### Important (Should Fix)

1. **No size/resource acceptance criterion**
   - Sections: Snapshot rules; Validation
   - Issue: “Too large” has no bounded threshold or test.
   - Impact: A later manual snapshot could noticeably increase binary size and initialization I/O without detection.
   - Recommended fix: Pin an expected file count and total byte budget (or a repository-consistent maximum), then test it.

#### Minor (Nice to Have)

1. **Copy-cost boundary is implicit**
   - Section: Architecture and data flow
   - Recommendation: State that the snapshot is copied only during `dotagen init`, retaining the existing one-pass copy behavior.

### Security

#### Important (Should Fix)

1. **Trusted prompt-content integrity is not verified**
   - Section: Snapshot rules
   - Issue: The spec accepts a mutable upstream branch without requiring a pinned revision, reviewed manifest, or content-integrity check.
   - Impact: A changed or compromised upstream source could introduce hostile instructions into initialized user skill directories.
   - Recommended fix: Pin the commit; add a committed manifest containing imported paths and SHA-256 digests; require a human review for executable content, remote prompt loading, credential/exfiltration instructions, and unsafe command execution.

#### Minor (Nice to Have)

1. **License evidence and attribution are underspecified**
   - Sections: Snapshot rules; Documentation
   - Recommendation: Include the upstream MIT license or a provenance document that records the license, copyright/notice, and intentional frontmatter change.

### Process & operations

#### Important (Should Fix)

1. **Manual maintenance workflow is missing**
   - Section: Snapshot rules
   - Issue: The spec chooses manual updates but provides no reproducible update procedure.
   - Impact: Upstream drift and future upgrades are difficult to assess or release safely.
   - Recommended fix: Define the pinned-source retrieval, diff/license/file-set review, named test commands, and documentation-count update steps.

#### Minor (Nice to Have)

1. **CI/artifact verification is implicit**
   - Section: Validation
   - Recommendation: Name the expected Go test commands and, if practical, add an `init` smoke test against the embedded artifact.

### Product

#### Important (Should Fix)

1. **Discoverability metadata conflicts with the documentation promise**
   - Sections: Scope; Snapshot rules; Documentation
   - Issue: The spec promises Nutlope/Frontend & UI cataloging but permits only a frontmatter `name` change. Dotagen discovery uses frontmatter category/vendor fields.
   - Impact: CLI/dashboard discovery would not match the documented category and vendor.
   - Recommended fix: Explicitly require normalized frontmatter metadata: `name`, `description`, `category: "Frontend & UI"`, `vendor: "nutlope"`, and license metadata where supported; validate it.

2. **Vietnamese catalog scope is ambiguous**
   - Section: Documentation
   - Issue: Adding one entry to a legacy-oriented catalog may reinforce inconsistent totals, prefixes, and scope.
   - Impact: Users may not understand whether it describes legacy agents/skills or the built-in library.
   - Recommended fix: Explicitly choose to label it as legacy with a scope note, or modernize its introduction and structure before adding Hallmark.

#### Minor (Nice to Have)

1. **Snapshot provenance is not user-visible**
   - Section: Snapshot rules
   - Recommendation: Make the pinned source revision and snapshot date easily visible in the provenance record.

---

## Strengths

- The Architecture and data flow section uses the existing recursive embedded-file and `runInit` behavior, avoiding new runtime abstractions.
- The Goal, Non-goals, and default `targets: []` establish a safe, unsurprising user workflow: availability after initialization without automatic activation.
- The scope deliberately excludes network fetches and website/demo assets, yielding deterministic builds and limited initialization work.
- Validation targets the appropriate integration boundaries: embed discovery, nested-file copying, frontmatter, and default configuration.

## Recommendations

Resolve the import contract as a pinned, auditable snapshot (provenance, manifest, license attribution, bounded size, and manual-update checklist). Then make the intended vendor/category metadata and Vietnamese catalog scope explicit. No deployment, HA, authentication, database, or frontend UX requirements are warranted for this self-contained CLI change.

## Resolution

The clear findings were applied to the spec: immutable provenance, a
digest-bearing manifest, a manual update checklist, metadata requirements,
complete path and size validation, named test commands, and the initialization
copy boundary. The catalog-scope decision was resolved with the requester:
`docs/CATALOG.vi.md` remains a legacy catalog and will receive a short scope
note that points to the README for the current built-in library. The remaining
minor recommendations are either incorporated by these changes or deferred as
implementation-level test/CI details.
