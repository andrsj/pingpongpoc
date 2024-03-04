package logging

import (
	"fmt"
	"log/slog"
)

func GetAvailableLevels() []string {
	return []string{"debug", "info", "warning", "error"}
}

func StringToLevel(logLevel string) (slog.Level, error) {
	stringToLevel := map[string]slog.Level{
		"debug":   slog.LevelDebug,
		"info":    slog.LevelInfo,
		"warning": slog.LevelWarn,
		"error":   slog.LevelError,
	}

	if level, exists := stringToLevel[logLevel]; exists {
		return level, nil
	}
	return slog.LevelInfo, fmt.Errorf("unknown level: %s", logLevel)
}
