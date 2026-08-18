package skillsrc

import (
	"strings"
	"testing"
)

func TestReviewSkillsUseBoundedDifferentialWorkflow(t *testing.T) {
	t.Parallel()

	for _, skill := range []string{
		"dotagent-enolalabs-plan-review",
		"dotagent-enolalabs-spec-review",
	} {
		skill := skill
		t.Run(skill, func(t *testing.T) {
			t.Parallel()

			raw, err := ReadSkillFile(skill + "/SKILL.md")
			if err != nil {
				t.Fatal(err)
			}
			content := string(raw)

			for _, required := range []string{
				"**Focused**",
				"**Differential**",
				"maximum two reviewers",
				"Do not silently fall back to a full review",
				"read-only by default",
			} {
				if !strings.Contains(content, required) {
					t.Errorf("missing bounded-review rule %q", required)
				}
			}

			for _, obsolete := range []string{
				"All 5 groups (recommended)",
				"Never skip the user choice step",
				"The full plan content",
				"The full spec content",
			} {
				if strings.Contains(content, obsolete) {
					t.Errorf("contains obsolete expensive-review rule %q", obsolete)
				}
			}
		})
	}
}

func TestPlanReviewerDoesNotRequirePreimplementedCode(t *testing.T) {
	t.Parallel()

	raw, err := ReadSkillFile("dotagent-enolalabs-plan-review/reviewers/technical.md")
	if err != nil {
		t.Fatal(err)
	}
	content := string(raw)

	if strings.Contains(content, "Missing test code is always Critical") {
		t.Fatal("technical reviewer still requires full test implementation")
	}
	if strings.Contains(content, "Every code step must contain actual, copy-pasteable code") {
		t.Fatal("technical reviewer still requires full implementation code")
	}
	if !strings.Contains(content, "copy-pasteable test bodies are not required") {
		t.Fatal("technical reviewer does not state the plan sufficiency boundary")
	}
}
