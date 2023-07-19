package zerolog

import (
	"io"

	"darvaza.org/slog"
	"darvaza.org/slog/handlers/filter"
	"darvaza.org/slog/handlers/zerolog"
)

// New creates a new slog.Logger from the Config.
func (c Config) New() slog.Logger {
	var log slog.Logger

	if c.Logger != nil {
		log = zerolog.New(c.Logger)
	} else {
		log = zerolog.New(NewZerolog(c.Writer))
	}

	return filter.New(log, c.Level)
}

// New creates a new [slog.Logger] wrapper for a
// [zerolog.Logger] using the given [io.Writer] and
// restricted to entries above the given [slog.LogLevel]
// threshold.
func New(w io.Writer, level slog.LogLevel) slog.Logger {
	var c Config

	c.SetLevel(level)
	c.SetWriter(w)
	return c.New()
}
