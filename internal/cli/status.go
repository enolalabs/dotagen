package cli

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/enolalabs/dotagen/v2/internal/agent"
	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/platform"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show sync status of all agents and targets",
	Long:  "Display the current status of each agent for each target: synced, out-of-date, missing, or not targeted.",
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

		agentsDir := filepath.Join(dotgenDir, "agents")
		agents, err := agent.ParseAgentsDir(agentsDir)
		if err != nil {
			return err
		}

		agentMap := make(map[string]agent.Agent)
		for _, a := range agents {
			agentMap[a.Name] = a
		}

		registry := platform.NewRegistry()

		fmt.Printf("\nAgents: %d defined\n", len(agents))
		fmt.Printf("Targets: %d configured\n\n", len(cfg.Targets))

		for _, ag := range agents {
			fmt.Printf("  %s:\n", ag.Name)
			for _, target := range cfg.Targets {
				resolved := cfg.ResolveTargets(ag.Name)
				targeted := false
				for _, t := range resolved {
					if t == target {
						targeted = true
						break
					}
				}

				if !targeted {
					fmt.Printf("    ✗ %-12s (not targeted)\n", target)
					continue
				}

				adapter, _ := registry.Get(target)
				symlinkPath := filepath.Join(projectDir, adapter.SymlinkPath(ag.Name))
				status := checkStatus(ag, target, adapter, symlinkPath, dotgenDir)
				switch status {
				case "synced":
					fmt.Printf("    ✓ %-12s (synced)\n", target)
				case "out-of-date":
					fmt.Printf("    ⚠ %-12s (out of date)\n", target)
				case "missing":
					fmt.Printf("    ✗ %-12s (missing)\n", target)
				case "broken":
					fmt.Printf("    💔 %-11s (broken symlink)\n", target)
				}
			}
			fmt.Println()
		}

		return nil
	},
}

func checkStatus(ag agent.Agent, target string, adapter platform.Adapter, symlinkPath string, dotgenDir string) string {
	info, err := os.Lstat(symlinkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "missing"
		}
		return "missing"
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return "missing"
	}

	linkTarget, err := os.Readlink(symlinkPath)
	if err != nil {
		return "broken"
	}

	resolvedTarget := linkTarget
	if !filepath.IsAbs(resolvedTarget) {
		resolvedTarget = filepath.Join(filepath.Dir(symlinkPath), resolvedTarget)
	}

	if _, err := os.Stat(resolvedTarget); err != nil {
		return "broken"
	}

	generatedContent, err := os.ReadFile(resolvedTarget)
	if err != nil {
		return "broken"
	}

	rendered, err := adapter.Render(ag)
	if err != nil {
		return "out-of-date"
	}

	genHash := sha256.Sum256(generatedContent)
	renHash := sha256.Sum256([]byte(rendered))

	if genHash != renHash {
		return "out-of-date"
	}

	_ = dotgenDir
	return "synced"
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
