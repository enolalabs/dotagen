package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/k0walski/dotagen/internal/agent"
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
	description := ""
	if desc, ok := ag.Frontmatter["description"]; ok && desc != "" {
		description = desc
	} else {
		lines := strings.Split(ag.Content, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				description = strings.TrimPrefix(trimmed, "# ")
				break
			}
		}
	}

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
	return filepath.Join(".cursor", "rules", agentName+".mdc")
}

func (a *CursorAdapter) EnsureDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, ".cursor", "rules"), 0o755)
}
