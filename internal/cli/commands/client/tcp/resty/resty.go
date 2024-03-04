package resty

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func NewRestyCmd(logger *slog.Logger) *cobra.Command {
	// restyCmd represents the unix command
	restyCmd := &cobra.Command{
		Use:   "resty",
		Short: "resty command",
		Long:  `resty command`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.Info("restyTCP")
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("restyTCPCommand")
		},
	}
	return restyCmd
}
