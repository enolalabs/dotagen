package cli

import (
	"fmt"

	"github.com/k0walski/dotagen/internal/web"
	"github.com/spf13/cobra"
)

var (
	servePort int
	serveOpen bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web dashboard server",
	Long:  "Start a local HTTP server with a web dashboard for managing agents visually.",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := web.NewServer(".", servePort)
		if err != nil {
			return err
		}
		fmt.Printf("✓ dotagen dashboard running at http://localhost:%d\n", servePort)
		return server.Start()
	},
}

func init() {
	serveCmd.Flags().IntVar(&servePort, "port", 7890, "Port to serve the dashboard on")
	serveCmd.Flags().BoolVar(&serveOpen, "open", true, "Open browser automatically")
	rootCmd.AddCommand(serveCmd)
}
