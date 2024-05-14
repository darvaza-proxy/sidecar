package reconnect

import "io"

type (
	// Reader is an alias of the standard [io.Reader] implementing `Read([]byte) (int, error)`
	Reader = io.Reader
	// Writer is an alias of the standard [io.Writer] implementing `Write([]byte) (int, error)`
	Writer = io.Writer
	// Closer is an alias of the standard [io.Closer] implementing `Close() error`
	Closer = io.Closer
)

// A Flusher implements the Flush() error interface to write whatever is left on the buffer.
type Flusher interface {
	Flush() error
}

// A WriteFlusher implements [io.Writer] and [Flusher].
type WriteFlusher interface {
	Writer
	Flusher
}

// A WriteCloseFlusher implements [io.Writer], [io.Closer] and [Flusher].
type WriteCloseFlusher interface {
	Writer
	Closer
	Flusher
}
