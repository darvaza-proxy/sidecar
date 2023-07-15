package zerolog

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type (
	// ConsoleWriter is an alias of [zerolog.ConsoleWriter].
	ConsoleWriter = zerolog.ConsoleWriter
)

// DefaultLogWriter returns a zerolog.ConsoleWriter set to
// use os.Stderr as io.Writer.
func DefaultLogWriter() zerolog.ConsoleWriter {
	return zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stderr
	})
}

// SetConsoleWriter binds the Config to a new zerolog.ConsoleWriter.
func (c *Config) SetConsoleWriter(options ...func(*zerolog.ConsoleWriter)) {
	w := zerolog.NewConsoleWriter(options...)

	c.SetWriter(w)
}

// SetWriter binds the Config to a zerolog.Logger using a given writer
// or the DefaultLogWriter if none specified.
func (c *Config) SetWriter(w io.Writer) {
	if w == nil {
		w = DefaultLogWriter()
	}

	c.Writer = w
	c.Logger = NewZerolog(w)
}

// SetLogger specifies the zerolog.Logger to use.
func (c *Config) SetLogger(zlog *zerolog.Logger) {
	c.Writer = nil
	c.Logger = zlog
}

// NewZerolog is a shortcut to zerolog.New() that allows us
// to bypass zerolog shadowing, and use DefaultLogWriter
// if none is specified.
func NewZerolog(w io.Writer) *zerolog.Logger {
	if w == nil {
		w = DefaultLogWriter()
	}

	zlog := zerolog.New(w)
	return &zlog
}
