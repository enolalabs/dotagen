package cli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "dotagen",
	Short: "Define sub-agents once, inject everywhere.",
	Long: `dotagen — A CLI tool that lets you define coding sub-agents
in markdown and inject them into multiple coding agent platforms.

Supported platforms: Antigravity, Claude Code, Gemini CLI, OpenCode`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(banner())
		cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print dotagen version",
	Long:  "Print the dotagen version, build info, and supported platforms.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dotagen v%s\n", version)
		fmt.Printf("  go:      %s\n", runtime.Version())
		fmt.Printf("  os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("  platforms: antigravity, claude-code, gemini-cli, opencode\n")
	},
}

func banner() string {
	return fmt.Sprintf(`
  ____          _       _
 |  _ \ ___  __| | __ _| |_ ___ _ __
 | | | / _ \/ _`+"`"+` |/ _`+"`"+` | __/ _ \ '__|
 | |_| |  __/ (_| | (_| | ||  __/ |
 |____/ \___|\__,_|\__,_|\__\___|_|
 
 Define sub-agents once, inject everywhere. (v%s)`, version)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
