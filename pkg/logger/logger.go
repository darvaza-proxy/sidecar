// Package logger provides a logging implementation based on uber-go/zap.
package logger

import (
	"darvaza.org/slog"
	"darvaza.org/slog/handlers/filter"
	"darvaza.org/slog/handlers/zap"
)

// New creates a darvaza.org/slog.Logger using the given
// go.uber.org/zap.Config.
// If cfg is nil, a default console logger will be created with
// development-friendly settings.
func New(cfg *Config) slog.Logger {
	if cfg == nil {
		cfg = zap.NewDefaultConfig()
	}
	return zap.New(cfg)
}

// NewWithThreshold works like [New] but restricts output to
// messages at or above the specified level threshold.
//
// Example thresholds:
//
//	slog.LevelDebug   // Show all messages
//	slog.LevelInfo    // Show info and above
//	slog.LevelWarning // Show only warnings and errors
func NewWithThreshold(cfg *Config, threshold slog.LogLevel) slog.Logger {
	return filter.New(New(cfg), threshold)
}

// NewDefaultConfig returns the console config for developers
// [New] will use when no config is provided.
func NewDefaultConfig() *Config {
	return zap.NewDefaultConfig()
}
