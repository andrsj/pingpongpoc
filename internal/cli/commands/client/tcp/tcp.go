package tcp

import (
	"log/slog"

	"pingpongpoc/internal/cli/commands/client/tcp/resty"
	"pingpongpoc/internal/cli/commands/client/tcp/std"

	"github.com/spf13/cobra"
)

func NewTCPCmd(logger *slog.Logger) *cobra.Command {
	// tcpCmd represents the tcp command
	tcpCmd := &cobra.Command{
		Use:   "tcp",
		Short: "tcp subcommands",
		Long:  `tcp subcommands`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.Info("tcp")
		},
		// Run: func(cmd *cobra.Command, args []string) {
		// 	logger.Info("TCPCommand")
		// },
	}

	tcpCmd.AddCommand(std.NewStdCmd(logger))
	tcpCmd.AddCommand(resty.NewRestyCmd(logger))

	return tcpCmd
}
