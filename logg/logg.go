package logg

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

type (
	FormatType int

	// LoggOpts customize [slog.HandlerOptions]
	LoggOpts struct {
		// Component enriches each log line with a componenent key/value.
		// Useful for aggregating/filtering with your log collector.
		Component string
		// Group nests individual keys in the format group.child.
		Group string
		// Log format.
		// Logfmt is the default log format
		// Human prints colourized logs useful for CLIs or development
		FormatType FormatType
		// Minimal level to log.
		// Debug level will automatically enable source code location.
		LogLevel slog.Level
	}
)

const (
	Logfmt FormatType = iota
	Human
	JSON
)

// NewLogg creates a new logger with the given options.
// Default log format is Logfmt.
func NewLogg(o LoggOpts) *slog.Logger {
	w := os.Stderr

	switch o.FormatType {
	case Human:
		handlerOpts := &tint.Options{
			Level:      o.LogLevel,
			TimeFormat: time.Kitchen,
		}
		if o.LogLevel == slog.LevelDebug {
			handlerOpts.AddSource = true
		}
		return enrichLogger(slog.New(tint.NewHandler(w, handlerOpts)), o.Component, o.Group)
	case Logfmt:
		return enrichLogger(slog.New(slog.NewTextHandler(w, populateHandlerOpts(o))), o.Component, o.Group)
	case JSON:
		return enrichLogger(slog.New(slog.NewJSONHandler(w, populateHandlerOpts(o))), o.Component, o.Group)
	default:
		return enrichLogger(slog.New(slog.NewTextHandler(w, populateHandlerOpts(o))), o.Component, o.Group)
	}
}

func populateHandlerOpts(o LoggOpts) *slog.HandlerOptions {
	handlerOpts := &slog.HandlerOptions{
		Level: o.LogLevel,
	}
	if o.LogLevel == slog.LevelDebug {
		handlerOpts.AddSource = true
	}
	return handlerOpts
}

func enrichLogger(baseLogger *slog.Logger, component string, group string) *slog.Logger {
	if component != "" {
		baseLogger = baseLogger.With("component", component)
	}
	if group != "" {
		baseLogger = baseLogger.WithGroup(group)
	}
	return baseLogger
}
