package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/enolalabs/dotagen/v2/internal/agent"
	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/skill"
)

// AntigravityAdapter implements the Adapter and SkillAdapter interfaces
// for Google Antigravity IDE.
//
// Antigravity uses:
//   - .agents/ directory at project root for workspace-scoped agents/skills
//   - SKILL.md files with YAML frontmatter (name, description)
//   - Plain markdown content for agent instructions (no frontmatter needed)
type AntigravityAdapter struct{}

func NewAntigravityAdapter() *AntigravityAdapter {
	return &AntigravityAdapter{}
}

func (a *AntigravityAdapter) Name() string {
	return "antigravity"
}

// Render outputs plain markdown content — Antigravity does not require
// YAML frontmatter for agent instructions (similar to Gemini CLI).
func (a *AntigravityAdapter) Render(ag agent.Agent) (string, error) {
	return ag.Content, nil
}

func (a *AntigravityAdapter) OutputPath(agentName string) string {
	return filepath.Join("antigravity", agentName+".md")
}

func (a *AntigravityAdapter) SymlinkPath(agentName string) string {
	return filepath.Join(config.AntigravityRootPath, agentName+".md")
}

func (a *AntigravityAdapter) EnsureDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.AntigravityRootPath), 0o755)
}

// SkillAdapter implementation — Antigravity uses SKILL.md with YAML
// frontmatter containing name and description.

func (a *AntigravityAdapter) RenderSkill(sk skill.Skill) (string, error) {
	desc := skill.ExtractDescription(sk)

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("name: %s\n", sk.Name))
	if desc != "" {
		sb.WriteString(fmt.Sprintf("description: %s\n", desc))
	}
	sb.WriteString("---\n\n")
	sb.WriteString(sk.Content)
	return sb.String(), nil
}

func (a *AntigravityAdapter) SkillOutputDir(skillName string) string {
	return filepath.Join("antigravity", "skills", skillName)
}

func (a *AntigravityAdapter) SkillSymlinkDir(skillName string) string {
	return filepath.Join(config.AntigravitySkillPath, skillName)
}

func (a *AntigravityAdapter) EnsureSkillDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.AntigravitySkillPath), 0o755)
}
