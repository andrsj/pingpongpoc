package std

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func NewStdCmd(logger *slog.Logger) *cobra.Command {
	// stdCmd represents the unix command
	stdCmd := &cobra.Command{
		Use:   "std",
		Short: "std subcommands",
		Long:  `std subcommands`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.Info("stdUnix")
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("stdUnixCommand")
		},
	}
	return stdCmd
}
