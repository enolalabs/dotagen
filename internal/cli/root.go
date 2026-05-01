package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "2.0.0"

var rootCmd = &cobra.Command{
	Use:   "dotagen",
	Short: "Define sub-agents once, inject everywhere.",
	Long: `dotagen — A CLI tool that lets you define coding sub-agents
in markdown and inject them into multiple coding agent platforms.

Supported platforms: Claude Code, Cursor, Gemini CLI, OpenCode`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(banner())
		cmd.Help()
	},
}

func banner() string {
	return `
  ____          _       _
 |  _ \ ___  __| | __ _| |_ ___ _ __
 | | | / _ \/ _` + "`" + ` |/ _` + "`" + ` | __/ _ \ '__|
 | |_| |  __/ (_| | (_| | ||  __/ |
 |____/ \___|\__,_|\__,_|\__\___|_|

 Define sub-agents once, inject everywhere.`
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
