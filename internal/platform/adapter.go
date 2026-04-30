package platform

import (
	"github.com/enolalabs/dotagen/internal/agent"
)

type Adapter interface {
	Name() string
	Render(agent agent.Agent) (string, error)
	OutputPath(agentName string) string
	SymlinkPath(agentName string) string
	EnsureDirectories(projectDir string) error
}
