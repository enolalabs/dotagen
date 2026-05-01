package platform

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/enolalabs/dotagen/v2/internal/agent"
	"github.com/enolalabs/dotagen/v2/internal/config"
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
