package platform

import (
	"os"
	"path/filepath"

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

func (a *GeminiCLIAdapter) Render(ag agent.Agent) (string, error) {
	return ag.Content, nil
}

func (a *GeminiCLIAdapter) OutputPath(agentName string) string {
	return filepath.Join("gemini-cli", agentName+".md")
}

func (a *GeminiCLIAdapter) SymlinkPath(agentName string) string {
	return filepath.Join(".gemini", "agents", agentName+".md")
}

func (a *GeminiCLIAdapter) EnsureDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, ".gemini", "agents"), 0o755)
}

// SkillAdapter implementation — Gemini CLI uses plain SKILL.md without transformation.

func (a *GeminiCLIAdapter) RenderSkill(sk skill.Skill) (string, error) {
	return sk.Content, nil
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
