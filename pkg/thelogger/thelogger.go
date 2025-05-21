package thelogger

import (
	"log/slog"
	"os"
)

type TheLogger struct {
	*slog.Logger // Embedded with this TheLogger inherits all slog.Logger methods
}

// NewTheLogger create a new instance of TheLogger
func NewTheLogger() *TheLogger {
	jsonHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true, // Add file and line
		Level:     slog.LevelDebug,
	})

	baseLogger := slog.New(jsonHandler)

	return &TheLogger{
		Logger: baseLogger,
	}
}
