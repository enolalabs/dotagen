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

// CodexAdapter implements Adapter and SkillAdapter for OpenAI Codex CLI.
//
// Codex uses:
//   - AGENTS.md at project root (or ~/.codex/AGENTS.md globally)
//   - Single-file format: each agent is rendered as a section within AGENTS.md
//   - For dotagen compatibility, we use .codex/agents/ directory with individual
//     agent files, and generate a combined AGENTS.md in the project root.
//
// Agent file path: .codex/agents/<name>.md  (individual agent files)
// Skill file path: .codex/skills/<name>/SKILL.md
type CodexAdapter struct{}

func NewCodexAdapter() *CodexAdapter {
	return &CodexAdapter{}
}

func (a *CodexAdapter) Name() string {
	return "codex"
}

// Render outputs markdown with YAML frontmatter for Codex.
// Each agent is rendered as its own .md file in .codex/agents/,
// matching the pattern used by other platforms.
func (a *CodexAdapter) Render(ag agent.Agent) (string, error) {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("name: %s\n", ag.Name))
	if desc, ok := ag.Frontmatter["description"]; ok && desc != "" {
		sb.WriteString(fmt.Sprintf("description: %s\n", desc))
	}
	sb.WriteString("---\n\n")
	sb.WriteString(ag.Content)
	return sb.String(), nil
}

func (a *CodexAdapter) OutputPath(agentName string) string {
	return filepath.Join("codex", agentName+".md")
}

func (a *CodexAdapter) SymlinkPath(agentName string) string {
	return filepath.Join(config.CodexRootPath, agentName+".md")
}

func (a *CodexAdapter) EnsureDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.CodexRootPath), 0o755)
}

// SkillAdapter implementation

func (a *CodexAdapter) RenderSkill(sk skill.Skill) (string, error) {
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

func (a *CodexAdapter) SkillOutputDir(skillName string) string {
	return filepath.Join("codex", "skills", skillName)
}

func (a *CodexAdapter) SkillSymlinkDir(skillName string) string {
	return filepath.Join(config.CodexSkillPath, skillName)
}

func (a *CodexAdapter) EnsureSkillDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.CodexSkillPath), 0o755)
}
