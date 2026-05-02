package cli

import (
	"fmt"
	"path/filepath"

	"github.com/enolalabs/dotagen/v2/internal/agent"
	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/engine"
	"github.com/enolalabs/dotagen/v2/internal/platform"
	"github.com/enolalabs/dotagen/v2/internal/skill"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync [target]",
	Short: "Sync agents and skills to target platforms",
	Long:  "Render all agents and skills, then create symlinks. Optionally specify a single target (claude-code, cursor, gemini-cli, opencode).",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dotgenDir, err := config.FindDotgenDir()
		if err != nil {
			return fmt.Errorf("failed to find dotgen directory: %w", err)
		}

		projectDir, err := config.GetProjectDir()
		if err != nil {
			return fmt.Errorf("failed to get project directory: %w", err)
		}

		cfg, err := config.LoadConfig(dotgenDir)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid config: %w", err)
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
			Skills:  make(map[string]config.SkillConfig),
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
		for name, sc := range cfg.Skills {
			resolved := cfg.ResolveSkillTargets(name)
			var filtered config.StringOrSlice
			for _, t := range resolved {
				for _, st := range syncTargets {
					if t == st {
						filtered = append(filtered, t)
					}
				}
			}
			if len(filtered) > 0 {
				filteredCfg.Skills[name] = config.SkillConfig{Targets: filtered, Disabled: sc.Disabled}
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
			return fmt.Errorf("failed to render agents: %w", err)
		}

		var agentNames []string
		for _, ag := range agents {
			agentNames = append(agentNames, ag.Name)
		}

		removed, err := engine.RemoveStaleSymlinks(projectDir, dotgenDir, agentNames, syncTargets)
		if err != nil {
			fmt.Printf("  ⚠ Failed to clean stale agent symlinks: %v\n", err)
		}

		// Parse and render skills
		skillsDir := filepath.Join(dotgenDir, "skills")
		skills, err := skill.ParseSkillsDir(skillsDir)
		if err != nil {
			fmt.Printf("  ⚠ Failed to parse skills: %v\n", err)
			skills = nil
		}

		var skillResults []engine.SkillRenderResult
		if len(skills) > 0 {
			skillResults, err = renderer.RenderAllSkills(skills, filteredCfg, dotgenDir, projectDir)
			if err != nil {
				fmt.Printf("  ⚠ Failed to render skills: %v\n", err)
			}
		}

		var skillNames []string
		for _, sk := range skills {
			skillNames = append(skillNames, sk.Name)
		}
		removedSkills, err := engine.RemoveStaleSkillSymlinks(projectDir, dotgenDir, skillNames, syncTargets)
		if err != nil {
			fmt.Printf("  ⚠ Failed to clean stale skill symlinks: %v\n", err)
		}

		totalSynced := len(results) + len(skillResults)
		fmt.Printf("✓ Synced %d agent(s) and %d skill(s) to %d platform(s)\n\n", len(results), len(skillResults), len(syncTargets))

		grouped := make(map[string][]engine.RenderResult)
		for _, r := range results {
			grouped[r.Target] = append(grouped[r.Target], r)
		}

		groupedSkills := make(map[string][]engine.SkillRenderResult)
		for _, r := range skillResults {
			groupedSkills[r.Target] = append(groupedSkills[r.Target], r)
		}

		for _, target := range syncTargets {
			fmt.Printf("  %s:\n", target)
			for _, r := range grouped[target] {
				relSymlink, _ := filepath.Rel(projectDir, r.SymlinkPath)
				relGenerated, _ := filepath.Rel(projectDir, r.GeneratedPath)
				fmt.Printf("    ✓ %s → %s\n", relSymlink, relGenerated)
			}
			for _, r := range groupedSkills[target] {
				relSymlink, _ := filepath.Rel(projectDir, r.SymlinkPath)
				relGenerated, _ := filepath.Rel(projectDir, r.GeneratedPath)
				fmt.Printf("    ✓ %s → %s\n", relSymlink, relGenerated)
			}
			fmt.Println()
		}

		allRemoved := append(removed, removedSkills...)
		if len(allRemoved) > 0 {
			fmt.Println("  Removed stale symlinks:")
			for _, r := range allRemoved {
				fmt.Printf("    ✗ %s\n", r)
			}
			fmt.Println()
		}

		_ = totalSynced

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
