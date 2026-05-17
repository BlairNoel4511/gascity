// Package logger provides a structured logging interface for gascity.
// It wraps the standard slog package with application-specific configuration.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Level represents the logging level.
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Logger wraps slog.Logger with additional context.
type Logger struct {
	*slog.Logger
}

// New creates a new Logger instance configured based on the provided level and format.
// format can be "json" or "text" (default).
func New(level Level, format string) *Logger {
	var logLevel slog.Level
	switch strings.ToLower(string(level)) {
	case string(LevelDebug):
		logLevel = slog.LevelDebug
	case string(LevelWarn):
		logLevel = slog.LevelWarn
	case string(LevelError):
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
		AddSource: logLevel == slog.LevelDebug,
	}

	var handler slog.Handler
	if strings.ToLower(format) == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// With returns a new Logger that includes the given attributes in every log record.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

// WithComponent returns a new Logger scoped to a specific application component.
func (l *Logger) WithComponent(component string) *Logger {
	return l.With("component", component)
}

// WithRequestID returns a new Logger that includes the given request ID.
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.With("request_id", requestID)
}
