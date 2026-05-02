package platform

import (
	"github.com/enolalabs/dotagen/v2/internal/skill"
)

// SkillAdapter defines how a platform renders and places skill files.
type SkillAdapter interface {
	Name() string
	RenderSkill(sk skill.Skill) (string, error)
	SkillOutputDir(skillName string) string   // e.g. "claude-code/skills/ds-tdd"
	SkillSymlinkDir(skillName string) string  // e.g. ".claude/skills/ds-tdd"
	EnsureSkillDirectories(projectDir string) error
}
