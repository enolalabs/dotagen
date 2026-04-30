package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/enolalabs/dotagen/internal/builtin"
	"github.com/enolalabs/dotagen/internal/engine"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize .dotagen directory structure",
	Long:  "Create a .dotagen/ directory with 144 built-in agents and default config. All agents are created with empty targets — edit config.yaml to enable them.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		dotgenDir := filepath.Join(cwd, ".dotagen")
		projectDir := cwd

		if _, err := os.Stat(dotgenDir); err == nil {
			fmt.Printf(".dotagen/ already exists at %s\n", dotgenDir)
			fmt.Print("Overwrite? (y/N): ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input != "y" && input != "yes" {
				fmt.Println("Aborted.")
				return nil
			}

			links, _ := engine.FindDotagenSymlinks(projectDir)
			for _, link := range links {
				os.Remove(link.Path)
			}
			engine.RemoveGeneratedContents(dotgenDir)
			os.RemoveAll(filepath.Join(dotgenDir, "agents"))
			os.Remove(filepath.Join(dotgenDir, "config.yaml"))
		}

		dirs := []string{
			filepath.Join(dotgenDir, "agents"),
			filepath.Join(dotgenDir, ".generated"),
		}
		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}

		agentNames := builtin.ListAgents()
		for _, name := range agentNames {
			data, err := builtin.ReadAgent(name)
			if err != nil {
				return fmt.Errorf("failed to read built-in agent %s: %w", name, err)
			}
			outPath := filepath.Join(dotgenDir, "agents", name+".md")
			if err := os.WriteFile(outPath, data, 0o644); err != nil {
				return fmt.Errorf("failed to write agent %s: %w", name, err)
			}
		}

		configContent := buildDefaultConfig(agentNames)
		if err := os.WriteFile(filepath.Join(dotgenDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
			return err
		}

		gitignoreContent := ".generated/\n"
		if err := os.WriteFile(filepath.Join(dotgenDir, ".gitignore"), []byte(gitignoreContent), 0o644); err != nil {
			return err
		}

		fmt.Printf("✓ Created .dotagen/ with %d built-in agents\n\n", len(agentNames))
		fmt.Println("  .dotagen/")
		fmt.Println("    config.yaml       — Configuration file (all agents disabled by default)")
		fmt.Println("    agents/           — Agent definitions (*.md)")
		fmt.Println("    .generated/       — Generated output (git-ignored)")
		fmt.Println("    .gitignore")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Edit .dotagen/config.yaml to set targets for agents you want to use")
		fmt.Println("  2. Run 'dotagen sync' to generate and symlink files")
		fmt.Println()
		fmt.Println("Example — enable da-backend-developer for all platforms:")
		fmt.Println("  agents:")
		fmt.Println("    da-backend-developer:")
		fmt.Println("      targets: all")

		return nil
	},
}

func buildDefaultConfig(agentNames []string) string {
	var sb strings.Builder
	sb.WriteString("# dotagen configuration\n")
	sb.WriteString("# Docs: https://github.com/enolalabs/dotagen\n")
	sb.WriteString("#\n")
	sb.WriteString("# All agents are listed with empty targets (disabled).\n")
	sb.WriteString("# Set targets to enable them. Examples:\n")
	sb.WriteString("#   targets: all                    — all platforms\n")
	sb.WriteString("#   targets: [claude-code, cursor]  — specific platforms\n")
	sb.WriteString("\n")
	sb.WriteString("targets:\n")
	sb.WriteString("  - claude-code\n")
	sb.WriteString("  - cursor\n")
	sb.WriteString("  - gemini-cli\n")
	sb.WriteString("  - opencode\n")
	sb.WriteString("\n")
	sb.WriteString("agents:\n")
	for _, name := range agentNames {
		sb.WriteString(fmt.Sprintf("  %s:\n    targets: []\n", name))
	}
	return sb.String()
}

func init() {
	rootCmd.AddCommand(initCmd)
}
