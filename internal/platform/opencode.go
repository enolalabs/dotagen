package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/enolalabs/dotagen/v2/internal/agent"
	"github.com/enolalabs/dotagen/v2/internal/config"
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
	return filepath.Join(config.OPEN_CODE_ROOT_PATH, agentName+".md")
}

func (a *OpenCodeAdapter) EnsureDirectories(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, config.OPEN_CODE_ROOT_PATH), 0o755)
}
