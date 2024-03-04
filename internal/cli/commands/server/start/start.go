package start

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func NewStartCmd(logger *slog.Logger) *cobra.Command {
	// startCmd represents the start command
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Starting server",
		Long:  `Starting server`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.Info("start")
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("StartCommand")
		},
	}

	return startCmd
}
