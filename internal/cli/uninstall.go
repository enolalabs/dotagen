package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/engine"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Completely remove dotagen from your system",
	Long: `Remove all dotagen artifacts from your system:

  1. All agent symlinks (da-*) in platform directories
  2. All skill symlinks (ds-*) in platform directories
  3. Antigravity global workflow files (~/.gemini/antigravity/global_workflows/da-*.md)
  4. The ~/.dotagen/ directory (config, agents, skills, generated files)
  5. The dotagen binary itself

This is a destructive operation and cannot be undone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		dotgenDir := filepath.Join(home, ".dotagen")
		if _, err := os.Stat(dotgenDir); os.IsNotExist(err) {
			fmt.Println("dotagen is not installed (~/.dotagen/ not found).")
			return nil
		}

		if !force {
			fmt.Println("⚠ This will completely remove dotagen from your system:")
			fmt.Println()
			fmt.Printf("  • All agent/skill symlinks in platform directories\n")
			fmt.Printf("  • Global workflow files in ~/.gemini/antigravity/global_workflows/\n")
			fmt.Printf("  • The entire ~/.dotagen/ directory\n")
			fmt.Printf("  • The dotagen binary\n")
			fmt.Println()
			fmt.Print("Are you sure? (yes/N): ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			input = strings.TrimSpace(strings.ToLower(input))
			if input != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		removed := 0

		// 1. Remove agent symlinks (da-*)
		links, err := engine.FindDotagenSymlinks(home, dotgenDir)
		if err != nil {
			fmt.Printf("  ⚠ Failed to find agent symlinks: %v\n", err)
		} else {
			for _, link := range links {
				if err := os.Remove(link.Path); err != nil {
					fmt.Printf("  ✗ Failed to remove %s: %v\n", link.Path, err)
					continue
				}
				rel, _ := filepath.Rel(home, link.Path)
				fmt.Printf("  ✓ Removed agent symlink ~/%s\n", rel)
				removed++
			}
		}

		// 2. Remove skill symlinks (ds-*)
		skillLinks, err := engine.FindDotagenSkillSymlinks(home, dotgenDir)
		if err != nil {
			fmt.Printf("  ⚠ Failed to find skill symlinks: %v\n", err)
		} else {
			for _, link := range skillLinks {
				if err := os.RemoveAll(link.Path); err != nil {
					fmt.Printf("  ✗ Failed to remove %s: %v\n", link.Path, err)
					continue
				}
				rel, _ := filepath.Rel(home, link.Path)
				fmt.Printf("  ✓ Removed skill symlink ~/%s\n", rel)
				removed++
			}
		}

		// 3. Remove Antigravity global workflow files (da-*.md)
		workflowsDir := filepath.Join(home, config.AntigravityGlobalWorkflowsPath)
		if entries, err := os.ReadDir(workflowsDir); err == nil {
			wfRemoved := 0
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				name := entry.Name()
				if strings.HasPrefix(name, "da-") && strings.HasSuffix(name, ".md") {
					fp := filepath.Join(workflowsDir, name)
					if err := os.Remove(fp); err != nil {
						fmt.Printf("  ✗ Failed to remove workflow %s: %v\n", name, err)
						continue
					}
					wfRemoved++
				}
			}
			if wfRemoved > 0 {
				fmt.Printf("  ✓ Removed %d global workflow file(s) from ~/%s\n", wfRemoved, config.AntigravityGlobalWorkflowsPath)
				removed += wfRemoved
			}
		}

		// 4. Remove empty platform directories left behind
		cleanEmptyPlatformDirs(home)

		// 5. Remove ~/.dotagen/ directory
		if err := os.RemoveAll(dotgenDir); err != nil {
			fmt.Printf("  ✗ Failed to remove ~/.dotagen/: %v\n", err)
		} else {
			fmt.Println("  ✓ Removed ~/.dotagen/")
			removed++
		}

		// 6. Remove the binary itself
		binaryPath, err := os.Executable()
		if err == nil {
			binaryPath, err = filepath.EvalSymlinks(binaryPath)
		}
		if err == nil {
			if err := os.Remove(binaryPath); err != nil {
				fmt.Printf("  ⚠ Could not remove binary %s: %v\n", binaryPath, err)
				fmt.Printf("    Remove it manually: rm %s\n", binaryPath)
			} else {
				fmt.Printf("  ✓ Removed binary %s\n", binaryPath)
				removed++
			}
		}

		fmt.Printf("\n✓ Uninstalled dotagen (%d item(s) removed)\n", removed)
		return nil
	},
}

// cleanEmptyPlatformDirs removes dotagen-created platform directories
// if they are empty after symlink removal.
func cleanEmptyPlatformDirs(homeDir string) {
	dirs := []string{
		config.ClaudeCodeRootPath,
		config.ClaudeCodeSkillPath,
		config.CodexRootPath,
		config.CodexSkillPath,
		config.GeminiCliRootPath,
		config.GeminiCliSkillPath,
		config.OpenCodeRootPath,
		config.OpenCodeSkillPath,
		config.AntigravityRootPath,
		config.AntigravitySkillPath,
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(homeDir, dir)
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			continue
		}
		if len(entries) == 0 {
			os.Remove(fullPath)
		}
	}
}

func init() {
	uninstallCmd.Flags().Bool("force", false, "Skip confirmation prompt")
	rootCmd.AddCommand(uninstallCmd)
}
