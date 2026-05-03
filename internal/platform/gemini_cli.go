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

type GeminiCLIAdapter struct{}

func NewGeminiCLIAdapter() *GeminiCLIAdapter {
	return &GeminiCLIAdapter{}
}

func (a *GeminiCLIAdapter) Name() string {
	return "gemini-cli"
}

// Render outputs markdown with Gemini CLI-standard YAML frontmatter
// containing name, description, tools, and model fields.
func (a *GeminiCLIAdapter) Render(ag agent.Agent) (string, error) {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("name: %s\n", ag.Name))
	if desc, ok := ag.Frontmatter["description"]; ok && desc != "" {
		sb.WriteString(fmt.Sprintf("description: %s\n", desc))
	}
	if tools, ok := ag.Frontmatter["tools"]; ok && tools != "" {
		sb.WriteString(fmt.Sprintf("tools:\n"))
		for _, tool := range strings.Split(tools, ",") {
			tool = strings.TrimSpace(tool)
			if tool != "" {
				sb.WriteString(fmt.Sprintf("  - %s\n", tool))
			}
		}
	}
	if model, ok := ag.Frontmatter["model"]; ok && model != "" {
		sb.WriteString(fmt.Sprintf("model: %s\n", model))
	}
	sb.WriteString("---\n\n")
	sb.WriteString(ag.Content)
	return sb.String(), nil
}

func (a *GeminiCLIAdapter) OutputPath(agentName string) string {
	return filepath.Join("gemini-cli", agentName+".md")
}

func (a *GeminiCLIAdapter) SymlinkPath(agentName string) string {
	return filepath.Join(config.GeminiCliRootPath, agentName+".md")
}

func (a *GeminiCLIAdapter) EnsureDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.GeminiCliRootPath), 0o755)
}

// SkillAdapter implementation — Gemini CLI uses SKILL.md with YAML
// frontmatter containing name and description.

func (a *GeminiCLIAdapter) RenderSkill(sk skill.Skill) (string, error) {
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

func (a *GeminiCLIAdapter) SkillOutputDir(skillName string) string {
	return filepath.Join("gemini-cli", "skills", skillName)
}

func (a *GeminiCLIAdapter) SkillSymlinkDir(skillName string) string {
	return filepath.Join(config.GeminiCliSkillPath, skillName)
}

func (a *GeminiCLIAdapter) EnsureSkillDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.GeminiCliSkillPath), 0o755)
}
