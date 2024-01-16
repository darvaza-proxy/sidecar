package zap

import (
	"darvaza.org/slog"
	"darvaza.org/slog/handlers/filter"
	"darvaza.org/slog/handlers/zap"
)

// New creates a new slog.Logger from the Config.
func (c Config) New() slog.Logger {
	log := zap.New(c.Config)
	return filter.New(log, c.Level)
}

// New creates a new [slog.Logger] wrapper for a
// [zap.Logger] using the given [io.Writer] and
// restricted to entries above the give [slog.LogLevel]
// threshold.
func New(level slog.LogLevel) slog.Logger {
	var c Config

	c.SetLevel(level)
	return c.New()
}
