package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/engine"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all generated files and symlinks",
	Long:  "Remove all symlinks created by dotagen and clean the .dotagen/.generated/ directory.",
	RunE: func(cmd *cobra.Command, args []string) error {
		dotgenDir, err := config.FindDotgenDir()
		if err != nil {
			return err
		}

		projectDir, err := config.GetProjectDir()
		if err != nil {
			return err
		}

		links, err := engine.FindDotagenSymlinks(projectDir, dotgenDir)
		if err != nil {
			return fmt.Errorf("failed to find symlinks: %w", err)
		}

		removed := 0
		for _, link := range links {
			if !strings.HasPrefix(link.Agent, "da-") {
				continue
			}
			if err := os.Remove(link.Path); err != nil {
				fmt.Printf("  ✗ Failed to remove %s: %v\n", link.Path, err)
				continue
			}
			rel, _ := filepath.Rel(projectDir, link.Path)
			fmt.Printf("  ✓ Removed %s\n", rel)
			removed++
		}

		if err := engine.RemoveGeneratedContents(dotgenDir); err != nil {
			return fmt.Errorf("failed to clean generated directory: %w", err)
		}

		fmt.Printf("\n✓ Cleaned %d symlink(s) and generated files\n", removed)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
