package server

import (
	"log/slog"

	"github.com/spf13/cobra"

	"pingpongpoc/internal/cli/commands/server/start"
)

func NewServerCmd(logger *slog.Logger) *cobra.Command {
	// ServerCmd represents the server command
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Server subcommands",
		Long:  `Server subcommands`,
		// PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// logger.Info("server")
		// },
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("ServerCommand")
		},
	}

	serverCmd.AddCommand(start.NewStartCmd(logger))

	return serverCmd
}
