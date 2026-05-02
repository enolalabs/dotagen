package platform

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/enolalabs/dotagen/v2/internal/agent"
	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/skill"
	"gopkg.in/yaml.v3"
)

type CursorFrontmatter struct {
	Description string `yaml:"description"`
	Globs       string `yaml:"globs,omitempty"`
	AlwaysApply bool   `yaml:"alwaysApply"`
}

type CursorAdapter struct{}

func NewCursorAdapter() *CursorAdapter {
	return &CursorAdapter{}
}

func (a *CursorAdapter) Name() string {
	return "cursor"
}

func (a *CursorAdapter) Render(ag agent.Agent) (string, error) {
	description := agent.ExtractDescription(ag)

	fm := CursorFrontmatter{
		Description: description,
		AlwaysApply: true,
	}

	fmBytes, err := yaml.Marshal(&fm)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor frontmatter: %w", err)
	}

	return fmt.Sprintf("---\n%s---\n\n%s", string(fmBytes), ag.Content), nil
}

func (a *CursorAdapter) OutputPath(agentName string) string {
	return filepath.Join("cursor", agentName+".mdc")
}

func (a *CursorAdapter) SymlinkPath(agentName string) string {
	return filepath.Join(config.CursorRootPath, agentName+".mdc")
}

func (a *CursorAdapter) EnsureDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.CursorRootPath), 0o755)
}

func (a *CursorAdapter) RenderSkill(sk skill.Skill) (string, error) {
	desc := skill.ExtractDescription(sk)
	fm := CursorFrontmatter{
		Description: desc,
		AlwaysApply: false,
	}
	fmBytes, err := yaml.Marshal(&fm)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor skill frontmatter: %w", err)
	}
	return fmt.Sprintf("---\n%s---\n\n%s", string(fmBytes), sk.Content), nil
}

func (a *CursorAdapter) SkillOutputDir(skillName string) string {
	return filepath.Join("cursor", "skills", skillName)
}

func (a *CursorAdapter) SkillSymlinkDir(skillName string) string {
	return filepath.Join(config.CursorSkillPath, skillName)
}

func (a *CursorAdapter) EnsureSkillDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.CursorSkillPath), 0o755)
}
