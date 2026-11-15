package logger

import (
	"io"
	"log/slog"
	"os"
)

// Logger wraps slog.Logger so both microservices use the same logging interface.
type Logger struct {
	*slog.Logger
}

// New creates a logger that writes structured JSON logs.
func New() *Logger {
	// Write logs to standard output so Docker and AWS can collect them.
	return NewWithOutput(os.Stdout)
}

// NewWithOutput allows us to choose the log destination.
// This is useful for testing and redirecting logs.
func NewWithOutput(output io.Writer) *Logger {
	// Configure JSON logging with a minimum level of INFO.
	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Create the underlying slog logger.
	baseLogger := slog.New(handler)

	// Return our reusable Logger wrapper.
	return &Logger{
		Logger: baseLogger,
	}
}
