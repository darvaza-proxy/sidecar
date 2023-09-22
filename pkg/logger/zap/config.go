// Package zap provides a zap based logger we use in our servers
package zap

import (
	"go.uber.org/zap"

	"darvaza.org/slog"
)

// Config describes how the Zap wrapper will be created
type Config struct {
	// Level is the threshold for the created slog.Logger
	Level slog.LogLevel

	Config *zap.Config
}

// SetLevel sets the threshold level for the created slog.Logger.
func (c *Config) SetLevel(level slog.LogLevel) {
	c.Level = level
}
