# Hallmark Built-in Skill Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox syntax for tracking.

**Goal:** Bundle the pinned Nutlope Hallmark design skill as a disabled-by-default Dotagent built-in skill with auditable provenance and complete initialization coverage.

**Architecture:** Add the static snapshot under skillsrc/data/dotagent-nutlope-hallmark. Existing go:embed and runInit discovery/copy behavior requires no production Go changes. SOURCE.md is a generated ledger with source identity, complete path/digest list, file count, and byte count.

**Tech Stack:** Go 1.26.2, embed.FS, github.com/adrg/frontmatter, SHA-256, Markdown.

## Global Constraints

- Import only skills/hallmark/SKILL.md and references/ from Nutlope/hallmark commit aeb42fb354ff4efa36ab475773a082315a3af2ce.
- The upstream snapshot is 106 Markdown files and 675085 bytes. SOURCE.md is Dotagent provenance, excluded from those totals.
- Retain upstream body and relative links. Normalize target metadata to name dotagent:nutlope:hallmark, vendor nutlope, category Frontend & UI, description, and MIT license.
- Do not alter scripts/fetch-official-skills.py, add network behavior, or enable the new skill.
- README must state 769 skills, 58 vendors, and 24 Frontend & UI skills. CATALOG.vi.md stays legacy and gets only a scope note.
- Required checks: go test ./skillsrc ./internal/cli and go test ./....

---

## File Structure

- Create: skillsrc/data/dotagent-nutlope-hallmark/SKILL.md
- Create: skillsrc/data/dotagent-nutlope-hallmark/references/** (105 upstream files)
- Create: skillsrc/data/dotagent-nutlope-hallmark/SOURCE.md
- Create: skillsrc/hallmark_snapshot_test.go
- Create: internal/cli/init_test.go
- Modify: README.md
- Modify: docs/CATALOG.vi.md

### Task 1: Define the failing import contract

**Files:**

- Create: skillsrc/hallmark_snapshot_test.go
- Create: internal/cli/init_test.go

**Interfaces:**

- Consumes: ListSkills(), ReadSkillFile(path string), ListSkillFiles(skillName string), and runInit(home string) (int, int, []string, error).
- Produces: test gates for the pinned snapshot, every embedded nested file, and targets: [].

- [ ] **Step 1: Create the failing embedded-snapshot test**

Write skillsrc/hallmark_snapshot_test.go. Parse SOURCE.md frontmatter into map[string]any, then assert source_repository, source_path, source_commit, license, upstream_files, and upstream_bytes. Convert files to []any, assert exactly 106 entries, then for every map entry read the stated path through ReadSkillFile, calculate sha256.Sum256, compare its lowercase hexadecimal digest, and accumulate len(raw) to equal 675085.

At the start of this test, also require that ListSkills() contains dotagent-nutlope-hallmark; this protects the embedded directory-discovery contract independently of direct file reads.

~~~go
func TestHallmarkSnapshotIsCompleteAndPinned(t *testing.T) {
	raw, err := ReadSkillFile("dotagent-nutlope-hallmark/SOURCE.md")
	if err != nil { t.Fatal(err) }

	var manifest map[string]any
	if _, err := frontmatter.Parse(strings.NewReader(string(raw)), &manifest); err != nil { t.Fatal(err) }
	if manifest["source_repository"] != "https://github.com/Nutlope/hallmark" ||
		manifest["source_path"] != "skills/hallmark" ||
		manifest["source_commit"] != "aeb42fb354ff4efa36ab475773a082315a3af2ce" ||
		manifest["license"] != "MIT" {
		t.Fatalf("unexpected provenance: %#v", manifest)
	}

	files, ok := manifest["files"].([]any)
	if !ok || len(files) != 106 { t.Fatalf("files = %#v", manifest["files"]) }
	total := 0
	for _, value := range files {
		file := value.(map[string]any)
		path, digest := file["path"].(string), file["sha256"].(string)
		if path == "SOURCE.md" || strings.HasPrefix(path, "/") || strings.Contains(path, "..") {
			t.Fatalf("invalid path %q", path)
		}
		content, err := ReadSkillFile("dotagent-nutlope-hallmark/" + path)
		if err != nil { t.Fatalf("read %s: %v", path, err) }
		if got := fmt.Sprintf("%x", sha256.Sum256(content)); got != digest {
			t.Fatalf("digest %s = %s, want %s", path, got, digest)
		}
		total += len(content)
	}
	if total != 675085 { t.Fatalf("bytes = %d, want 675085", total) }
}
~~~

Add TestHallmarkSkillMetadata in the same file. Parse SKILL.md frontmatter to map[string]string and require name dotagent:nutlope:hallmark, vendor nutlope, category Frontend & UI, license MIT, and a non-empty description.

- [ ] **Step 2: Create the failing runInit test**

Write internal/cli/init_test.go in package cli:

~~~go
func TestRunInitCopiesHallmarkAndLeavesItDisabled(t *testing.T) {
	home := t.TempDir()
	if _, _, _, err := runInit(home); err != nil { t.Fatal(err) }

	const name = "dotagent-nutlope-hallmark"
	for _, rel := range skillsrc.ListSkillFiles(name) {
		want, err := skillsrc.ReadSkillFile(name + "/" + rel)
		if err != nil { t.Fatalf("embedded %s: %v", rel, err) }
		got, err := os.ReadFile(filepath.Join(home, ".dotagen", "skills", name, rel))
		if err != nil { t.Fatalf("initialized %s: %v", rel, err) }
		if !bytes.Equal(got, want) { t.Fatalf("%s differs from embedded content", rel) }
	}
	config, err := os.ReadFile(filepath.Join(home, ".dotagen", "config.yaml"))
	if err != nil { t.Fatal(err) }
	if !strings.Contains(string(config), "  "+name+":\n    targets: []\n") {
		t.Fatalf("Hallmark is not disabled: %s", config)
	}
}
~~~

Imports are bytes, os, path/filepath, strings, testing, and github.com/enolalabs/dotagen/v2/skillsrc.

- [ ] **Step 3: Verify red state**

Run: go test ./skillsrc ./internal/cli -run 'Test(Hallmark|RunInitCopiesHallmark)' -count=1

Expected: FAIL because the embedded Hallmark directory does not exist.

- [ ] **Step 4: Commit the contract**

~~~bash
git add skillsrc/hallmark_snapshot_test.go internal/cli/init_test.go
git commit -m "test: define Hallmark built-in skill contract"
~~~

### Task 2: Import the pinned snapshot and provenance ledger

**Files:**

- Create: skillsrc/data/dotagent-nutlope-hallmark/SKILL.md
- Create: skillsrc/data/dotagent-nutlope-hallmark/references/**
- Create: skillsrc/data/dotagent-nutlope-hallmark/SOURCE.md

**Interfaces:**

- Consumes: Task 1 test contract and the fixed upstream SHA.
- Produces: the complete directory consumed by existing embed and runInit behavior; no production Go source changes.

- [ ] **Step 1: Checkout and verify the exact upstream revision**

~~~bash
git clone https://github.com/Nutlope/hallmark.git /tmp/hallmark-import
git -C /tmp/hallmark-import checkout --detach aeb42fb354ff4efa36ab475773a082315a3af2ce
git -C /tmp/hallmark-import rev-parse HEAD
find /tmp/hallmark-import/skills/hallmark -type f | wc -l
du -sb /tmp/hallmark-import/skills/hallmark
~~~

Expected: the fixed SHA, 106 files, and 675085 bytes. Stop if any differs.

- [ ] **Step 2: Review import safety and license**

~~~bash
find /tmp/hallmark-import/skills/hallmark -type f ! -name '*.md' -print
rg -n -i 'curl |wget |fetch\(|https?://|api[_-]?key|token|secret|password|rm -rf|sudo ' /tmp/hallmark-import/skills/hallmark
sed -n '1,24p' /tmp/hallmark-import/LICENSE
~~~

Expected: only Markdown is imported. Review every search hit as static text guidance; reject remote prompt loading, credential exfiltration, or unsafe command execution. Confirm the MIT notice is Copyright (c) 2026 Hallmark contributors.

- [ ] **Step 3: Copy only the allowed directory and normalize metadata**

~~~bash
mkdir -p skillsrc/data/dotagent-nutlope-hallmark
cp -R /tmp/hallmark-import/skills/hallmark/. skillsrc/data/dotagent-nutlope-hallmark/
find skillsrc/data/dotagent-nutlope-hallmark -type f | wc -l
~~~

Expected: 106 files. Replace only target SKILL.md frontmatter with:

~~~yaml
---
name: "dotagent:nutlope:hallmark"
description: "Anti-AI-slop design skill for greenfield pages, audits, redesigns, and design extraction from URLs or screenshots."
version: "1.1.0"
category: "Frontend & UI"
vendor: "nutlope"
license: "MIT"
---
~~~

Keep everything after the frontmatter unchanged.

- [ ] **Step 4: Generate SOURCE.md**

Use YAML frontmatter matching this shape, then add a concise manual-update checklist: checkout the pin, inspect license/files/diff/text safety, regenerate digests, update README counts, run focused and full Go tests.

~~~yaml
---
source_repository: "https://github.com/Nutlope/hallmark"
source_path: "skills/hallmark"
source_commit: "aeb42fb354ff4efa36ab475773a082315a3af2ce"
snapshot_date: "2026-07-22"
license: "MIT"
copyright: "Copyright (c) 2026 Hallmark contributors"
upstream_files: 106
upstream_bytes: 675085
files:
  - path: "SKILL.md"
    sha256: "<post-normalization digest>"
---
~~~

Populate files with all 106 sorted upstream paths and their SHA-256 values after metadata normalization; exclude SOURCE.md itself.

~~~bash
find skillsrc/data/dotagent-nutlope-hallmark -type f ! -name SOURCE.md -printf '%P\n' | sort | while read -r path; do sha256sum "skillsrc/data/dotagent-nutlope-hallmark/$path"; done
~~~

- [ ] **Step 5: Verify green state and commit**

Run: go test ./skillsrc ./internal/cli -run 'Test(Hallmark|RunInitCopiesHallmark)' -count=1

Expected: PASS; the tests prove every digest, byte count, nested-file copy, and disabled config entry.

~~~bash
git add skillsrc/data/dotagent-nutlope-hallmark skillsrc/hallmark_snapshot_test.go internal/cli/init_test.go
git commit -m "feat: add Hallmark built-in skill"
~~~

### Task 3: Reconcile public documentation

**Files:**

- Modify: README.md
- Modify: docs/CATALOG.vi.md

**Interfaces:**

- Consumes: the imported snapshot and final directory count.
- Produces: accurate library statistics and a clear legacy documentation boundary.

- [ ] **Step 1: Establish stale documentation**

Run: rg -n '740 official skills|55 vendors|\| Frontend & UI \| 23 \||all 740 built-in skills|768 skills' README.md

Expected: matches existing stale claims.

- [ ] **Step 2: Apply exact README changes**

Replace each built-in library 740 count with 769; change 55 vendors and View all 55 vendors to 58; change Frontend & UI 23 to 24; change Skill Library matrix 768 to 769. Add:

~~~markdown
| Nutlope | 1 | Frontend & UI |
~~~

Run: rg -n '740|55 vendors|\| Frontend & UI \| 23 \||768 skills' README.md

Expected: no matches.

- [ ] **Step 3: Preserve the Vietnamese catalog as legacy**

Immediately after the title in docs/CATALOG.vi.md add:

~~~markdown
> **Phạm vi lịch sử:** Tài liệu này mô tả catalog agents/skills `ds-*` trước đây. Danh sách built-in skills hiện hành, số lượng và cách bật `dotagent-*` được duy trì trong [README](../README.md#built-in-skills).
~~~

Do not add Hallmark to legacy tables or alter their historical ds-* entries.

- [ ] **Step 4: Verify documentation and commit**

~~~bash
test "$(find skillsrc/data -mindepth 1 -maxdepth 1 -type d -name 'dotagent-*' | wc -l | tr -d ' ')" = 769
git diff --check
git add README.md docs/CATALOG.vi.md
git commit -m "docs: catalog Hallmark built-in skill"
~~~

Expected: both checks pass.

### Task 4: Run release-quality verification

**Files:** none.

**Interfaces:** Consumes Tasks 1–3 and produces a review-ready working tree.

- [ ] **Step 1: Run focused tests**

Run: go test ./skillsrc ./internal/cli

Expected: PASS, including all Hallmark contract tests.

- [ ] **Step 2: Run full suite**

Run: go test ./...

Expected: PASS.

- [ ] **Step 3: Verify change boundary**

~~~bash
git status --short
git diff --check HEAD
git log -3 --oneline
~~~

Expected: feature commits contain only snapshot, tests, README, and legacy catalog work; unrelated work is untouched.

- [ ] **Step 4: Record handoff**

Report the fixed SHA, 106-file/675085-byte totals, both test commands, and PASS results.

## Self-Review

- **Spec coverage:** Tasks 1–2 cover exact provenance, metadata, digest integrity, file/size guard, manual-update safety, and disabled initialization. Task 3 covers the current library docs and agreed legacy boundary. Task 4 runs both required suites.
- **Placeholder scan:** The digest example is generated by the stated command and must be expanded to every imported file; no implementation decision is deferred.
- **Type consistency:** Task 2 emits the manifest keys and file paths consumed by Task 1; Task 1 calls existing skillsrc and runInit signatures unchanged.


## Plan Review Amendments

These amendments are part of the tasks above and supersede any less-specific wording.

### Task 1 amendments: executable exact-set contract

Add this complete metadata test to skillsrc/hallmark_snapshot_test.go:

~~~go
func TestHallmarkSkillMetadata(t *testing.T) {
	raw, err := ReadSkillFile("dotagent-nutlope-hallmark/SKILL.md")
	if err != nil { t.Fatal(err) }
	var metadata map[string]string
	if _, err := frontmatter.Parse(strings.NewReader(string(raw)), &metadata); err != nil { t.Fatal(err) }
	for key, want := range map[string]string{
		"name": "dotagent:nutlope:hallmark",
		"vendor": "nutlope",
		"category": "Frontend & UI",
		"license": "MIT",
	} {
		if got := metadata[key]; got != want { t.Errorf("%s = %q, want %q", key, got, want) }
	}
	if metadata["description"] == "" { t.Error("description is empty") }
}
~~~

In TestHallmarkSnapshotIsCompleteAndPinned, decode upstream_files and upstream_bytes from the manifest as integers; fail if they are missing/non-integer, if len(files) differs from upstream_files, or if the summed bytes differ from upstream_bytes. Then require upstream_files == 106 and upstream_bytes == 675085.

Build a sorted manifest-path set with checked map/string assertions. Reject duplicates, SOURCE.md, absolute paths, parent segments, non-md paths, and paths other than SKILL.md or references/*.md. Build a sorted embedded-path set from ListSkillFiles("dotagent-nutlope-hallmark"), exclude SOURCE.md, and require deep equality with the manifest set before hashing. In the init test, sort ListSkillFiles(name), read every listed destination, and separately assert that no destination files exist outside this list.

### Task 2 amendments: safe, deterministic import

Before copying, execute these rejection checks:

~~~bash
test -z "$(find /tmp/hallmark-import/skills/hallmark -type l -print -quit)"
test -z "$(find /tmp/hallmark-import/skills/hallmark ! -type f -print -quit)"
test -z "$(find /tmp/hallmark-import/skills/hallmark -type f ! -name '*.md' -print -quit)"
~~~

Expected: all commands exit 0. Fail the import on any symlink, non-regular entry, or non-Markdown file.

Use the regex scan only as triage. Also review the complete pinned tree with:

~~~bash
git -C /tmp/hallmark-import diff --no-index -- /dev/null skills/hallmark
rg -n '/(home|tmp|Users)/|file://' /tmp/hallmark-import/skills/hallmark
~~~

Record this sign-off in SOURCE.md prose: reviewer examined the complete file set and confirmed no remote prompt loading, credential collection/exfiltration, executable payloads, unsafe command instructions, or absolute filesystem references that break post-init use.

Generate SOURCE.md deterministically with this temporary Go program, run from the repository root after SKILL.md metadata has been normalized:

~~~go
package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	root := "skillsrc/data/dotagent-nutlope-hallmark"
	var paths []string
	var total int
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if info.IsDir() { return nil }
		rel, err := filepath.Rel(root, path)
		if err != nil { return err }
		if rel != "SOURCE.md" { paths = append(paths, filepath.ToSlash(rel)) }
		return nil
	})
	if err != nil { panic(err) }
	sort.Strings(paths)
	fmt.Println("---")
	fmt.Println("source_repository: \"https://github.com/Nutlope/hallmark\"")
	fmt.Println("source_path: \"skills/hallmark\"")
	fmt.Println("source_commit: \"aeb42fb354ff4efa36ab475773a082315a3af2ce\"")
	fmt.Println("snapshot_date: \"2026-07-22\"")
	fmt.Println("license: \"MIT\"")
	fmt.Println("copyright: \"Copyright (c) 2026 Hallmark contributors\"")
	fmt.Printf("upstream_files: %d\n", len(paths))
	fmt.Println("upstream_bytes: 675085")
	fmt.Println("files:")
	for _, rel := range paths {
		raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil { panic(err) }
		total += len(raw)
		fmt.Printf("  - path: %q\n    sha256: \"%x\"\n", rel, sha256.Sum256(raw))
	}
	if total != 675085 || len(paths) != 106 { panic(fmt.Sprintf("files=%d bytes=%d", len(paths), total)) }
	fmt.Println("---")
	fmt.Println()
	fmt.Println("# Hallmark source provenance")
	fmt.Println()
	fmt.Println("Manual update: checkout the pinned commit; review the full file set and license; reject symlinks, non-Markdown files, unsafe instructions, and absolute filesystem references; regenerate this ledger; update README; run the required Go tests.")
}
~~~

Run it exactly as:

~~~bash
go run /tmp/hallmark-manifest.go > skillsrc/data/dotagent-nutlope-hallmark/SOURCE.md
rm -f /tmp/hallmark-manifest.go
rm -rf /tmp/hallmark-import
~~~

The generator writes a sorted full ledger; its fixed totals make failure explicit.

### Task 3 amendments: accurate provenance documentation

In README Built-in Skills introduction, replace the source sentence with:

~~~markdown
dotagen ships with **769 official skills** from **58 vendors**. Most are sourced from the [awesome-agent-skills](https://github.com/enolalabs/awesome-agent-skills) registry and [obra/superpowers](https://github.com/obra/superpowers); [Hallmark](https://github.com/Nutlope/hallmark) is a separately pinned MIT snapshot from Nutlope. They are injected automatically when you run `dotagen init`.
~~~

Apply the same distinction in the final Built-in Skills acknowledgment paragraph, retaining its registry/vendor wording while adding the linked Nutlope Hallmark MIT snapshot sentence.

### Task 4 amendments: durable CI gate and acceptance verification

Create .github/workflows/test.yml:

~~~yaml
name: Test

on:
  push:
  pull_request:

permissions:
  contents: read

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v6
        with:
          go-version: "1.26.2"
      - run: go test ./...
~~~

Before committing, run:

~~~bash
test -f .github/workflows/test.yml
rg -n 'go test ./\.\.\.' .github/workflows/test.yml
git add .github/workflows/test.yml
git commit -m "ci: verify Go test suite"
~~~

Add an optional handoff acceptance scenario: initialize with a temporary HOME, set dotagent-nutlope-hallmark targets to [codex] in config.yaml, run the existing sync/list command, and confirm the rendered Codex skill is present. When reviewing git status, identify only task-owned paths; do not require a clean worktree because unrelated user changes must be preserved.
