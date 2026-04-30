package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/k0walski/dotagen/internal/agent"
	"github.com/k0walski/dotagen/internal/config"
	"github.com/k0walski/dotagen/internal/platform"
)

type Renderer struct {
	registry *platform.Registry
}

func NewRenderer(registry *platform.Registry) *Renderer {
	return &Renderer{registry: registry}
}

func (r *Renderer) Render(ag agent.Agent, target string) (string, error) {
	adapter, err := r.registry.Get(target)
	if err != nil {
		return "", err
	}
	return adapter.Render(ag)
}

type RenderResult struct {
	AgentName   string
	Target      string
	GeneratedPath string
	SymlinkPath   string
}

func (r *Renderer) RenderAll(agents []agent.Agent, cfg *config.Config, dotgenDir string, projectDir string) ([]RenderResult, error) {
	var results []RenderResult
	generatedDir := filepath.Join(dotgenDir, ".generated")

	for _, ag := range agents {
		targets := cfg.ResolveTargets(ag.Name)
		for _, target := range targets {
			adapter, err := r.registry.Get(target)
			if err != nil {
				return nil, err
			}

			rendered, err := adapter.Render(ag)
			if err != nil {
				return nil, fmt.Errorf("failed to render agent %q for %s: %w", ag.Name, target, err)
			}

			outPath := filepath.Join(generatedDir, adapter.OutputPath(ag.Name))
			if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
				return nil, fmt.Errorf("failed to create output directory: %w", err)
			}

			if err := os.WriteFile(outPath, []byte(rendered), 0o644); err != nil {
				return nil, fmt.Errorf("failed to write generated file: %w", err)
			}

			absGenerated, err := filepath.Abs(outPath)
			if err != nil {
				return nil, err
			}

			symlinkPath := filepath.Join(projectDir, adapter.SymlinkPath(ag.Name))
			if err := adapter.EnsureDirectories(projectDir); err != nil {
				return nil, fmt.Errorf("failed to ensure directories for %s: %w", target, err)
			}

			if err := CreateSymlink(absGenerated, symlinkPath); err != nil {
				return nil, fmt.Errorf("failed to create symlink %s: %w", symlinkPath, err)
			}

			results = append(results, RenderResult{
				AgentName:     ag.Name,
				Target:        target,
				GeneratedPath: outPath,
				SymlinkPath:   symlinkPath,
			})
		}
	}

	return results, nil
}
