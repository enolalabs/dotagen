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

type ClaudeCodeAdapter struct{}

func NewClaudeCodeAdapter() *ClaudeCodeAdapter {
	return &ClaudeCodeAdapter{}
}

func (a *ClaudeCodeAdapter) Name() string {
	return "claude-code"
}

func (a *ClaudeCodeAdapter) Render(ag agent.Agent) (string, error) {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("name: %s\n", ag.Name))
	if desc, ok := ag.Frontmatter["description"]; ok && desc != "" {
		sb.WriteString(fmt.Sprintf("description: %s\n", desc))
	}
	if model, ok := ag.Frontmatter["model"]; ok && model != "" {
		sb.WriteString(fmt.Sprintf("model: %s\n", model))
	}
	sb.WriteString("---\n\n")
	sb.WriteString(ag.Content)
	return sb.String(), nil
}

func (a *ClaudeCodeAdapter) OutputPath(agentName string) string {
	return filepath.Join("claude-code", agentName+".md")
}

func (a *ClaudeCodeAdapter) SymlinkPath(agentName string) string {
	return filepath.Join(".claude", "agents", agentName+".md")
}

func (a *ClaudeCodeAdapter) EnsureDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, ".claude", "agents"), 0o755)
}

func (a *ClaudeCodeAdapter) RenderSkill(sk skill.Skill) (string, error) {
	var sb strings.Builder
	sb.WriteString("---\n")
	name := strings.TrimPrefix(sk.Name, "ds-")
	sb.WriteString(fmt.Sprintf("name: %s\n", name))
	if desc, ok := sk.Frontmatter["description"]; ok && desc != "" {
		sb.WriteString(fmt.Sprintf("description: %s\n", desc))
	}
	sb.WriteString("---\n\n")
	sb.WriteString(sk.Content)
	return sb.String(), nil
}

func (a *ClaudeCodeAdapter) SkillOutputDir(skillName string) string {
	return filepath.Join("claude-code", "skills", skillName)
}

func (a *ClaudeCodeAdapter) SkillSymlinkDir(skillName string) string {
	return filepath.Join(config.ClaudeCodeSkillPath, skillName)
}

func (a *ClaudeCodeAdapter) EnsureSkillDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.ClaudeCodeSkillPath), 0o755)
}
