package main

import (
	"context"

	"github.com/spf13/pflag"

	"darvaza.org/sidecar/pkg/logger"
	"darvaza.org/sidecar/pkg/service"
	"darvaza.org/slog"
	"darvaza.org/slog/handlers/filter"
)

func getLogLevel(flags *pflag.FlagSet) slog.LogLevel {
	level := slog.Error

	if flags != nil {
		verbosity, err := flags.GetCount(verbosityFlag)
		if err == nil {
			level += slog.LogLevel(verbosity)
			switch {
			case level < slog.Error:
				level = slog.Error
			case level > slog.Debug:
				level = slog.Debug
			}
		}
	}

	return level
}

func getSystemLogger(ctx context.Context, flags *pflag.FlagSet) slog.Logger {
	svc, _ := service.GetService(ctx)
	if svc != nil {
		if !svc.Interactive() || WantsSyslog(flags) {
			return svc.SystemLogger()
		}
	}
	return nil
}

func getLogger(ctx context.Context, flags *pflag.FlagSet) slog.Logger {
	level := getLogLevel(flags)
	if log := getSystemLogger(ctx, flags); log != nil {
		return filter.New(log, level)
	}
	return newLoggerLevel(level)
}

func newLogger(flags *pflag.FlagSet) slog.Logger {
	return newLoggerLevel(getLogLevel(flags))
}

func newLoggerLevel(level slog.LogLevel) slog.Logger {
	return logger.NewWithThreshold(nil, level)
}

const (
	verbosityFlag      = "verbose"
	verbosityShortFlag = "v"
)

func init() {
	pFlags := rootCmd.PersistentFlags()
	pFlags.CountP(verbosityFlag, verbosityShortFlag, "increase verbosity")
}
