package platform

import (
	"os"
	"path/filepath"

	"github.com/enolalabs/dotagen/internal/agent"
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
