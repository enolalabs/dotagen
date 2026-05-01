package cli

import (
	"fmt"

	"github.com/enolalabs/dotagen/v2/internal/web"
	"github.com/spf13/cobra"
)

var (
	servePort int
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
	rootCmd.AddCommand(serveCmd)
}
