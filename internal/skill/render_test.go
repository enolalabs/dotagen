package skill

import (
	"strings"
	"testing"
)

func TestReferenceNames(t *testing.T) {
	sk := Skill{
		Name:    "ds-tdd",
		Content: "# TDD",
		References: []Reference{
			{Name: "tests.md"},
			{Name: "mocking.md"},
			{Name: "scripts/run.sh"},
		},
	}
	names := sk.ReferenceNames()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	if names[0] != "tests.md" || names[1] != "mocking.md" || names[2] != "scripts/run.sh" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestReferenceNames_Empty(t *testing.T) {
	sk := Skill{Name: "ds-simple", Content: "# Simple"}
	names := sk.ReferenceNames()
	if len(names) != 0 {
		t.Fatalf("expected 0 names, got %d", len(names))
	}
}

func TestContentWithReferences_NoRefs(t *testing.T) {
	sk := Skill{Name: "test", Content: "Hello world"}
	got := sk.ContentWithReferences("")
	if got != "Hello world" {
		t.Errorf("expected unchanged content, got %q", got)
	}
}

func TestContentWithReferences_MarkdownRef(t *testing.T) {
	sk := Skill{
		Name:    "ds-tdd",
		Content: "# TDD\n\nSee [tests.md](tests.md)",
		References: []Reference{
			{Name: "tests.md", Content: "# Test Guidelines\n\nWrite good tests."},
		},
	}
	got := sk.ContentWithReferences(".claude/skills/ds-tdd")

	if !strings.Contains(got, "## Bundled Reference Files") {
		t.Error("missing bundled reference files section")
	}
	if !strings.Contains(got, "### tests.md") {
		t.Error("missing reference heading for tests.md")
	}
	if !strings.Contains(got, "# Test Guidelines") {
		t.Error("markdown reference content not inlined")
	}
	// Markdown refs should NOT have a file path annotation
	if strings.Contains(got, "**File path**") {
		t.Error("markdown reference should not have file path annotation")
	}
}

func TestContentWithReferences_CodeRef(t *testing.T) {
	sk := Skill{
		Name:    "ds-git-guardrails",
		Content: "# Git Guardrails\n\nCopy [scripts/block.sh](scripts/block.sh)",
		References: []Reference{
			{Name: "scripts/block.sh", Content: "#!/bin/bash\necho 'BLOCKED'\nexit 2"},
		},
	}
	got := sk.ContentWithReferences(".claude/skills/ds-git-guardrails")

	// Code refs MUST have file path annotation
	if !strings.Contains(got, "**File path**: `.claude/skills/ds-git-guardrails/scripts/block.sh`") {
		t.Error("missing file path annotation for code reference")
	}
	// Code refs must be in a code fence
	if !strings.Contains(got, "```bash") {
		t.Error("shell script not wrapped in bash code fence")
	}
}

func TestContentWithReferences_EmptySymlinkDir(t *testing.T) {
	sk := Skill{
		Name:    "test",
		Content: "# Test",
		References: []Reference{
			{Name: "run.py", Content: "print('hello')"},
		},
	}
	got := sk.ContentWithReferences("")

	// With empty symlinkDir, code files should still be inlined but without path annotation
	if strings.Contains(got, "**File path**") {
		t.Error("should not have file path when symlinkDir is empty")
	}
	if !strings.Contains(got, "```python") {
		t.Error("python file not wrapped in python code fence")
	}
}
