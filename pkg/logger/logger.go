package logger

import (
	"context"
	"log/slog"
	"os"
)

// InitLogger initializes the logger with the specified format.
// If format is "simple", it logs only the message; otherwise, it uses JSON format.
func InitLogger(format string) {
	var handler slog.Handler

	if format == "simple" {
		handler = &simpleHandler{}
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true, // Include source file and line number
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// simpleHandler is a custom slog.Handler that logs only the message.
type simpleHandler struct{}

// Enabled reports whether the handler is enabled for the given level.
func (h *simpleHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// Handle logs only the message part of the record.
func (h *simpleHandler) Handle(_ context.Context, r slog.Record) error {
	_, err := os.Stdout.WriteString(r.Message + "\n")
	return err
}

// WithAttrs returns a new handler with the given attributes.
func (h *simpleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h // Attributes are ignored in this simple handler
}

// WithGroup returns a new handler with the given group name.
func (h *simpleHandler) WithGroup(name string) slog.Handler {
	return h // Groups are ignored in this simple handler
}
