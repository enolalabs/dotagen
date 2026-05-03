package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/enolalabs/dotagen/v2/internal/agent"
	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/spf13/cobra"
)

var (
	createDescription string
	createTargets     string
	createContent     string
	createFile        string
	createTemplate    bool
)

var agentNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new agent",
	Long: `Create a new agent definition file and register it in config.yaml.

Content can be provided via:
  --content "markdown..."   Inline content
  --file template.md        Read from file
  (no flag)                 Opens $EDITOR for interactive editing

Examples:
  dotagen create review-code
  dotagen create review-code -d "Senior code reviewer" -t claude-code,gemini-cli
  dotagen create my-agent --content "# My Agent\n\nInstructions..."
  dotagen create my-agent --file ./template.md --template`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if !agentNameRegex.MatchString(name) {
			return fmt.Errorf("invalid agent name %q: must start with alphanumeric and contain only letters, digits, hyphens, and underscores", name)
		}

		dotgenDir, err := config.FindDotgenDir()
		if err != nil {
			return err
		}

		agentsDir := filepath.Join(dotgenDir, "agents")
		agentPath := filepath.Join(agentsDir, name+".md")

		if _, err := os.Stat(agentPath); err == nil {
			return fmt.Errorf("agent %q already exists at %s", name, agentPath)
		}

		content, err := resolveContent(name)
		if err != nil {
			return err
		}

		if createDescription != "" {
			content = fmt.Sprintf("---\ndescription: %s\n---\n\n%s", createDescription, content)
		}

		if err := os.MkdirAll(agentsDir, 0o755); err != nil {
			return fmt.Errorf("failed to create agents directory: %w", err)
		}

		if err := os.WriteFile(agentPath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("failed to write agent file: %w", err)
		}

		targets := parseTargetsFlag(createTargets)
		if err := config.AddAgentToConfig(dotgenDir, name, targets); err != nil {
			os.Remove(agentPath)
			return fmt.Errorf("failed to update config: %w", err)
		}

		relPath, _ := filepath.Rel(dotgenDir, agentPath)
		fmt.Printf("✓ Created agent %q\n", name)
		fmt.Printf("  File: %s\n", relPath)
		if len(targets) == 1 && targets[0] == "all" {
			fmt.Println("  Targets: all")
		} else {
			fmt.Printf("  Targets: %s\n", strings.Join(targets, ", "))
		}
		fmt.Println()
		fmt.Println("  Run 'dotagen sync' to generate platform files")

		return nil
	},
}

func resolveContent(name string) (string, error) {
	if createContent != "" {
		return createContent, nil
	}

	if createFile != "" {
		data, err := os.ReadFile(createFile)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", createFile, err)
		}
		return string(data), nil
	}

	if createTemplate {
		return agent.ScaffoldContent(name), nil
	}

	return interactiveContent(name)
}

func interactiveContent(name string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		return readFromStdin(name)
	}

	tmpFile, err := os.CreateTemp("", "dotagen-"+name+"-*.md")
	if err != nil {
		return readFromStdin(name)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	scaffold := agent.ScaffoldContent(name)
	tmpFile.WriteString(scaffold)
	tmpFile.Close()

	fmt.Printf("Opening editor for agent %q...\n", name)

	parts := strings.Fields(editor)
	args := append(parts[1:], tmpPath)
	execPath, err := findExec(parts[0])
	if err != nil {
		return readFromStdin(name)
	}

	exitCode := runEditor(execPath, args)
	if exitCode != 0 {
		return "", fmt.Errorf("editor exited with code %d", exitCode)
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to read editor output: %w", err)
	}

	content := strings.TrimSpace(string(data))
	if content == "" || content == strings.TrimSpace(scaffold) {
		return "", fmt.Errorf("empty agent content — agent not created")
	}

	return content, nil
}

func findExec(name string) (string, error) {
	return exec.LookPath(name)
}

func runEditor(execPath string, args []string) int {
	cmd := exec.Command(execPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return cmd.ProcessState.ExitCode()
}

func readFromStdin(name string) (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat stdin: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		content := strings.TrimSpace(string(data))
		if content == "" {
			return "", fmt.Errorf("empty input from stdin")
		}
		return content, nil
	}

	fmt.Printf("Enter agent content (Ctrl+D to finish):\n\n")
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	content := strings.TrimSpace(string(data))
	if content == "" {
		return "", fmt.Errorf("empty input")
	}
	return content, nil
}

func parseTargetsFlag(targets string) []string {
	if targets == "" {
		return []string{"all"}
	}
	parts := strings.Split(targets, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return []string{"all"}
	}
	return result
}

func init() {
	createCmd.Flags().StringVarP(&createDescription, "description", "d", "", "Agent description (frontmatter)")
	createCmd.Flags().StringVarP(&createTargets, "targets", "t", "all", "Target platforms (comma-separated)")
	createCmd.Flags().StringVarP(&createContent, "content", "c", "", "Inline markdown content")
	createCmd.Flags().StringVarP(&createFile, "file", "f", "", "Read content from file")
	createCmd.Flags().BoolVar(&createTemplate, "template", false, "Use built-in template scaffold")
	rootCmd.AddCommand(createCmd)
}
