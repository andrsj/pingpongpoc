package client

import (
	"log/slog"

	"pingpongpoc/internal/cli/commands/client/tcp"
	"pingpongpoc/internal/cli/commands/client/unix"

	"github.com/spf13/cobra"
)

func NewClientCmd(logger *slog.Logger) *cobra.Command {
	// clientCmd represents the client command
	clientCmd := &cobra.Command{
		Use:   "client",
		Short: "Client subcommands",
		Long:  `Client subcommands`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.Info("client")
		},
		// Run: func(cmd *cobra.Command, args []string) {
		// 	logger.Info("ClientCommand")
		// },
	}

	clientCmd.AddCommand(tcp.NewTCPCmd(logger))
	clientCmd.AddCommand(unix.NewUnixCmd(logger))

	return clientCmd
}
