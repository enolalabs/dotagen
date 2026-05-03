package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/enolalabs/dotagen/v2/internal/agent"
	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/skill"
	"gopkg.in/yaml.v3"
)

type OpenCodeFrontmatter struct {
	Description string `yaml:"description"`
	Mode        string `yaml:"mode"`
}

type OpenCodeAdapter struct{}

func NewOpenCodeAdapter() *OpenCodeAdapter {
	return &OpenCodeAdapter{}
}

func (a *OpenCodeAdapter) Name() string {
	return "opencode"
}

func (a *OpenCodeAdapter) Render(ag agent.Agent) (string, error) {
	description := agent.ExtractDescription(ag)

	fm := OpenCodeFrontmatter{
		Description: description,
		Mode:        "subagent",
	}

	fmBytes, err := yaml.Marshal(&fm)
	if err != nil {
		return "", fmt.Errorf("failed to marshal opencode frontmatter: %w", err)
	}

	return fmt.Sprintf("---\n%s---\n\n%s", string(fmBytes), ag.Content), nil
}

func (a *OpenCodeAdapter) OutputPath(agentName string) string {
	return filepath.Join("opencode", agentName+".md")
}

func (a *OpenCodeAdapter) SymlinkPath(agentName string) string {
	return filepath.Join(config.OpenCodeRootPath, agentName+".md")
}

func (a *OpenCodeAdapter) EnsureDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.OpenCodeRootPath), 0o755)
}

// SkillAdapter implementation

func (a *OpenCodeAdapter) RenderSkill(sk skill.Skill) (string, error) {
	desc := skill.ExtractDescription(sk)
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("name: %s\n", sk.Name))
	sb.WriteString(fmt.Sprintf("description: %s\n", desc))
	sb.WriteString("---\n\n")
	sb.WriteString(sk.Content)
	return sb.String(), nil
}

func (a *OpenCodeAdapter) SkillOutputDir(skillName string) string {
	return filepath.Join("opencode", "skills", skillName)
}

func (a *OpenCodeAdapter) SkillSymlinkDir(skillName string) string {
	return filepath.Join(config.OpenCodeSkillPath, skillName)
}

func (a *OpenCodeAdapter) EnsureSkillDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.OpenCodeSkillPath), 0o755)
}
