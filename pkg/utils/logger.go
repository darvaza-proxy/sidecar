package utils

import (
	"net"

	"darvaza.org/core"
	"darvaza.org/slog"
)

const (
	// LogFieldAddress is the field name used to store the address
	// when logging.
	LogFieldAddress = "addr"

	// LogFieldError is the field name used to store the error
	// when logging.
	LogFieldError = slog.ErrorFieldName
)

// LogWithAddress adds a [net.Addr] as a string field to the logger.
// it will do nothing if the logger
func LogWithAddress(l slog.Logger, label string, addr net.Addr) slog.Logger {
	if l != nil && !core.IsZero(addr) {
		if s := addr.String(); s != "" {
			if label == "" {
				label = LogFieldAddress
			}
			l = l.WithField(label, s)
		}
	}
	return l
}

// LogWithError adds the given error as field to the logger.
func LogWithError(l slog.Logger, err error) slog.Logger {
	if l != nil && err != nil {
		l = l.WithField(slog.ErrorFieldName, err)
	}
	return l
}
