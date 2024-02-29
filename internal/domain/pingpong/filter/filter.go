package filter

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Filter struct {
	log *slog.Logger
}

func New(logger *slog.Logger) *Filter {
	return &Filter{
		log: logger,
	}
}

func (f *Filter) Filter(ctx context.Context, input string) (bool, error) {
	f.log.Info("Filter executes", "input", input)

	select {
	case <-time.After(6 * time.Second): //nolint:gomnd // Simulate work by waiting for 3 seconds
	case <-ctx.Done():
		// If the context is canceled before the work completes, log and exit
		f.log.Info("Filter execution canceled", "input", input)

		return false, fmt.Errorf("%w", ctx.Err())
	}

	// Your existing input handling logic
	switch input {
	case "ping":
		return true, nil
	default:
		return false, nil
	}
}
