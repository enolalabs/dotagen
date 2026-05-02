package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/skill"
	"github.com/spf13/cobra"
)

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Manage skills (slash commands)",
	Long:  "Create, list, and delete skills. Skills are directory-based and contain SKILL.md + optional reference files.",
}

var skillListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all skills",
	RunE: func(cmd *cobra.Command, args []string) error {
		dotgenDir, err := config.FindDotgenDir()
		if err != nil {
			return err
		}

		cfg, err := config.LoadConfig(dotgenDir)
		if err != nil {
			return err
		}

		skills, err := skill.ParseSkillsDir(filepath.Join(dotgenDir, "skills"))
		if err != nil {
			return err
		}

		if len(skills) == 0 {
			fmt.Println("No skills found. Run 'dotagen init' to import built-in skills.")
			return nil
		}

		sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })

		fmt.Printf("\n%-35s %-12s %-15s %s\n", "SKILL", "CATEGORY", "TARGETS", "REFS")
		fmt.Println(strings.Repeat("─", 80))

		for _, sk := range skills {
			cat := sk.Frontmatter["category"]
			if cat == "" {
				cat = "—"
			}

			targets := cfg.ResolveSkillTargets(sk.Name)
			targetStr := "disabled"
			if len(targets) > 0 {
				targetStr = strings.Join(targets, ", ")
			}

			refCount := len(sk.References)
			refStr := "—"
			if refCount > 0 {
				refStr = fmt.Sprintf("%d file(s)", refCount)
			}

			fmt.Printf("%-35s %-12s %-15s %s\n", sk.Name, truncStr(cat, 12), truncStr(targetStr, 15), refStr)
		}

		fmt.Printf("\n%d skill(s) total\n", len(skills))
		return nil
	},
}

var skillCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new skill",
	Long:  "Create a new skill directory with a scaffold SKILL.md file.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dotgenDir, err := config.FindDotgenDir()
		if err != nil {
			return err
		}

		name := args[0]
		if !strings.HasPrefix(name, "ds-") {
			name = "ds-" + name
		}

		clean := strings.TrimPrefix(name, "ds-")
		if !isValidSkillCLIName(clean) {
			return fmt.Errorf("invalid skill name %q — use only lowercase letters, numbers, and hyphens", name)
		}

		skillDir := filepath.Join(dotgenDir, "skills", name)
		if _, err := os.Stat(skillDir); err == nil {
			return fmt.Errorf("skill %q already exists at %s", name, skillDir)
		}

		if err := os.MkdirAll(skillDir, 0o755); err != nil {
			return fmt.Errorf("failed to create skill directory: %w", err)
		}

		content := skill.ScaffoldSkillContent(name)
		skillFile := filepath.Join(skillDir, "SKILL.md")
		if err := os.WriteFile(skillFile, []byte(content), 0o644); err != nil {
			return fmt.Errorf("failed to write SKILL.md: %w", err)
		}

		// Add to config with all targets
		targets, _ := cmd.Flags().GetStringSlice("targets")
		if len(targets) == 0 {
			targets = []string{"all"}
		}
		if err := config.AddSkillToConfig(dotgenDir, name, targets); err != nil {
			fmt.Printf("  ⚠ Failed to add to config: %v\n", err)
		}

		fmt.Printf("✓ Created skill %s\n", name)
		fmt.Printf("  → %s\n\n", skillFile)
		fmt.Println("Next steps:")
		fmt.Println("  1. Edit the SKILL.md file with your skill workflow")
		fmt.Println("  2. Optionally add reference files to the skill directory")
		fmt.Println("  3. Run 'dotagen sync' to deploy")
		return nil
	},
}

var skillDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a skill",
	Long:  "Delete a skill directory and remove it from the config.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dotgenDir, err := config.FindDotgenDir()
		if err != nil {
			return err
		}

		name := args[0]
		if !strings.HasPrefix(name, "ds-") {
			name = "ds-" + name
		}

		skillDir := filepath.Join(dotgenDir, "skills", name)
		if _, err := os.Stat(skillDir); os.IsNotExist(err) {
			return fmt.Errorf("skill %q not found", name)
		}

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Printf("Delete skill %q? This cannot be undone. (use --force to skip)\n", name)
			return nil
		}

		if err := os.RemoveAll(skillDir); err != nil {
			return fmt.Errorf("failed to delete skill: %w", err)
		}

		if err := config.RemoveSkillFromConfig(dotgenDir, name); err != nil {
			fmt.Printf("  ⚠ Failed to remove from config: %v\n", err)
		}

		fmt.Printf("✓ Deleted skill %s\n", name)
		return nil
	},
}

func truncStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func isValidSkillCLIName(name string) bool {
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return len(name) > 0
}

func init() {
	skillCreateCmd.Flags().StringSlice("targets", nil, "Target platforms (default: all)")
	skillDeleteCmd.Flags().Bool("force", false, "Skip confirmation")

	skillCmd.AddCommand(skillListCmd)
	skillCmd.AddCommand(skillCreateCmd)
	skillCmd.AddCommand(skillDeleteCmd)
	rootCmd.AddCommand(skillCmd)
}
