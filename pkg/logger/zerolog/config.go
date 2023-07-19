// Package zerolog provides a zerolog based logger we use in our servers
package zerolog

import (
	"io"

	"darvaza.org/slog"
	"github.com/rs/zerolog"
)

// Config describes how the Zerolog wrapper will be created.
type Config struct {
	// Level is the threshold for the created slog.Logger
	Level slog.LogLevel

	// Logger optionally defines a zerolog.Logger to use
	Logger *zerolog.Logger

	// Writer defines where the zerolog.Logger would write,
	// if no Logger is specified
	Writer io.Writer
}

// SetLevel sets the threshold level for the created slog.Logger.
func (c *Config) SetLevel(level slog.LogLevel) {
	c.Level = level
}
