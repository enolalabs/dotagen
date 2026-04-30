package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
)

func ParseAgentsDir(agentsDir string) ([]Agent, error) {
	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read agents directory: %w", err)
	}

	var agents []Agent
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		filePath := filepath.Join(agentsDir, entry.Name())
		agent, err := ParseAgentFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse agent file %s: %w", entry.Name(), err)
		}
		agents = append(agents, *agent)
	}

	return agents, nil
}

func ParseAgentFile(filePath string) (*Agent, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	name := strings.TrimSuffix(filepath.Base(filePath), ".md")

	var fm map[string]string
	body, err := frontmatter.Parse(strings.NewReader(string(data)), &fm)
	if err != nil {
		body = data
		fm = make(map[string]string)
	}

	content := strings.TrimSpace(string(body))

	return &Agent{
		Name:        name,
		Content:     content,
		Frontmatter: fm,
		FilePath:    filePath,
	}, nil
}
