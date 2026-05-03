package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/enolalabs/dotagen/v2/internal/agent"
	"github.com/enolalabs/dotagen/v2/internal/config"
)

// SyncGlobalWorkflows writes all agents as workflow files to the
// Antigravity global_workflows directory. Each file uses the format:
//
//	---
//	description: <agent description>
//	---
//
//	<agent content>
//
// This is called automatically when syncing to the "antigravity" target.
func SyncGlobalWorkflows(agents []agent.Agent, projectDir string) (int, error) {
	workflowsDir := filepath.Join(projectDir, config.AntigravityGlobalWorkflowsPath)
	if err := os.MkdirAll(workflowsDir, 0o755); err != nil {
		return 0, fmt.Errorf("failed to create global_workflows directory: %w", err)
	}

	synced := 0
	for _, ag := range agents {
		desc := ag.Frontmatter["description"]

		var sb strings.Builder
		sb.WriteString("---\n")
		if desc != "" {
			sb.WriteString(fmt.Sprintf("description: %s\n", desc))
		} else {
			sb.WriteString(fmt.Sprintf("description: Use this agent: %s\n", ag.Name))
		}
		sb.WriteString("---\n\n")
		sb.WriteString(ag.Content)

		outPath := filepath.Join(workflowsDir, ag.Name+".md")
		if err := os.WriteFile(outPath, []byte(sb.String()), 0o644); err != nil {
			return synced, fmt.Errorf("failed to write workflow file %s: %w", ag.Name, err)
		}
		synced++
	}

	// Clean up stale workflow files that no longer have a matching agent
	activeSet := make(map[string]bool)
	for _, ag := range agents {
		activeSet[ag.Name+".md"] = true
	}

	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		return synced, nil // non-fatal — we synced what we could
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !activeSet[entry.Name()] {
			os.Remove(filepath.Join(workflowsDir, entry.Name()))
		}
	}

	return synced, nil
}
