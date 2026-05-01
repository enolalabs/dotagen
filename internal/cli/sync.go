package cli

import (
	"fmt"
	"path/filepath"

	"github.com/enolalabs/dotagen/internal/agent"
	"github.com/enolalabs/dotagen/internal/config"
	"github.com/enolalabs/dotagen/internal/engine"
	"github.com/enolalabs/dotagen/internal/platform"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync [target]",
	Short: "Sync agents to target platforms",
	Long:  "Render all agents and create symlinks. Optionally specify a single target (claude-code, cursor, gemini-cli, opencode).",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dotgenDir, err := config.FindDotgenDir()
		if err != nil {
			return err
		}

		projectDir, err := config.GetProjectDir()
		if err != nil {
			return err
		}

		cfg, err := config.LoadConfig(dotgenDir)
		if err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}

		var syncTargets []string
		if len(args) == 1 {
			syncTargets = []string{args[0]}
			valid := false
			for _, t := range config.ValidTargets {
				if t == args[0] {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid target %q; valid targets: %v", args[0], config.ValidTargets)
			}
		} else {
			detected := config.DetectPlatforms(projectDir)
			if len(detected) > 0 {
				syncTargets = detected
			} else {
				syncTargets = cfg.Targets
			}
		}

		filteredCfg := &config.Config{
			Targets: syncTargets,
			Agents:  make(map[string]config.AgentConfig),
		}
		for name, ac := range cfg.Agents {
			resolved := cfg.ResolveTargets(name)
			var filtered config.StringOrSlice
			for _, t := range resolved {
				for _, st := range syncTargets {
					if t == st {
						filtered = append(filtered, t)
					}
				}
			}
			if len(filtered) > 0 {
				filteredCfg.Agents[name] = config.AgentConfig{Targets: filtered, Disabled: ac.Disabled}
			}
		}

		agentsDir := filepath.Join(dotgenDir, "agents")
		agents, err := agent.ParseAgentsDir(agentsDir)
		if err != nil {
			return fmt.Errorf("failed to parse agents: %w", err)
		}

		if len(agents) == 0 {
			fmt.Println("No agents found in .dotagen/agents/")
			return nil
		}

		registry := platform.NewRegistry()
		renderer := engine.NewRenderer(registry)

		results, err := renderer.RenderAll(agents, filteredCfg, dotgenDir, projectDir)
		if err != nil {
			return err
		}

		var agentNames []string
		for _, ag := range agents {
			agentNames = append(agentNames, ag.Name)
		}

		removed, err := engine.RemoveStaleSymlinks(projectDir, agentNames, syncTargets)
		if err != nil {
			fmt.Printf("  ⚠ Failed to clean stale symlinks: %v\n", err)
		}

		fmt.Printf("✓ Synced %d agent(s) to %d platform(s)\n\n", len(results), len(syncTargets))

		grouped := make(map[string][]engine.RenderResult)
		for _, r := range results {
			grouped[r.Target] = append(grouped[r.Target], r)
		}

		for _, target := range syncTargets {
			fmt.Printf("  %s:\n", target)
			for _, r := range grouped[target] {
				relSymlink, _ := filepath.Rel(projectDir, r.SymlinkPath)
				relGenerated, _ := filepath.Rel(projectDir, r.GeneratedPath)
				fmt.Printf("    ✓ %s → %s\n", relSymlink, relGenerated)
			}
			fmt.Println()
		}

		if len(removed) > 0 {
			fmt.Println("  Removed stale symlinks:")
			for _, r := range removed {
				fmt.Printf("    ✗ %s\n", r)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
