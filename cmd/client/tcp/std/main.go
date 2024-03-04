package main

import (
	"context"
	"log/slog"
	"os"

	"pingpongpoc/internal/client/tcp/std"
	"pingpongpoc/internal/constants"
)

func main() {
	//nolint:exhaustruct
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}

	// Text handler
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	// JSON handler
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	client := std.NewClient(logger)

	ctx, cancel := context.WithTimeout(context.Background(), constants.TimeoutForPinging)
	defer cancel()

	err := client.Ping(ctx)
	if err != nil {
		logger.Error("Error", "err", err)
	}
}