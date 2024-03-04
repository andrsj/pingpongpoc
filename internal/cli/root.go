package cli

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"pingpongpoc/internal/cli/commands/client"
	"pingpongpoc/internal/cli/commands/server"
	"pingpongpoc/internal/logging"
)

const logLevelFlag = "log-level"

func NewPingPongCLI(logger *slog.Logger) *cobra.Command {
	// RootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "pingpongpoc",
		Short: "pingpong proof of concept application",
		Long: `Ping Pong (proof of concept) application
Built for testing different client/server architectures`,

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logLevelStr, err := cmd.Flags().GetString(logLevelFlag)
			if err != nil {
				logger.Error("Failed getting level for logging", "err", err)
			}

			fmt.Println("BLYAAAAAAAAAAAAAAAAAAAAAAAAA", logLevelStr)

			slog.SetLogLoggerLevel(slog.LevelError)

			if logLevelStr != "" {
				level, err := logging.StringToLevel(logLevelStr)
				if err != nil {
					logger.Error("Wrong logging level", "level", logLevelStr)
				}

				slog.SetLogLoggerLevel(level)

			}

			logger.Info("pingpongpoc")
		},
		// PreRun: func(cmd *cobra.Command, args []string) { fmt.Println("HUY") },
		// Run:    func(cmd *cobra.Command, args []string) { fmt.Println("HUY2") },
	}

	rootCmd.AddCommand(server.NewServerCmd(logger))
	rootCmd.AddCommand(client.NewClientCmd(logger))

	rootCmd.PersistentFlags().String(
		logLevelFlag,
		"error",
		fmt.Sprintf(
			"Messages with this level and above will be logged. Valid levels are: %s",
			strings.Join(logging.GetAvailableLevels(), ", "),
		),
	)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}

	// Text handler
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	// JSON handler
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	command := NewPingPongCLI(logger)

	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
