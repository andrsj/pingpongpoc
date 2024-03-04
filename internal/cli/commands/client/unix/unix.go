package unix

import (
	"log/slog"

	"pingpongpoc/internal/cli/commands/client/unix/resty"
	"pingpongpoc/internal/cli/commands/client/unix/std"

	"github.com/spf13/cobra"
)

func NewUnixCmd(logger *slog.Logger) *cobra.Command {
	// unixCmd represents the unix command
	unixCmd := &cobra.Command{
		Use:   "unix",
		Short: "Unix subcommands",
		Long:  `Unix subcommands`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.Info("unix")
		},
		// Run: func(cmd *cobra.Command, args []string) {
		// 	logger.Info("UnixCommand")
		// },
	}

	unixCmd.AddCommand(std.NewStdCmd(logger))
	unixCmd.AddCommand(resty.NewRestyCmd(logger))

	return unixCmd
}
