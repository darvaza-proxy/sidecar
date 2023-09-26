package service

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/spf13/pflag"
)

var (
	// ErrNotImplemented indicates something isn't implemented yet
	ErrNotImplemented = errors.New("not implemented")
)

const (
	ExitStatusOK    int = iota // ExitStatusOK indicates everything is fine
	ExitStatusMinor            // ExitStatusMinor indicates there was a minor problem
	ExitStatusMajor            // ExitStatusMajor indicates there was a major problem
)

// ErrorExitCode is an error wrapper that knows how the application
// should exit
type ErrorExitCode struct {
	Code int
	Err  error
}

func (e ErrorExitCode) Unwrap() error {
	return e.Err
}

// ExitStatus tells the code for os.Exit()
func (e ErrorExitCode) ExitStatus() int {
	return e.Code
}

func (e ErrorExitCode) Error() string {
	var buf bytes.Buffer

	if e.Code == 0 {
		return "OK"
	}

	_, _ = fmt.Fprint(&buf, "ExitCode=", e.Code)
	if e.Err != nil {
		_, _ = fmt.Fprint(&buf, ": ", e.Err.Error())
	}

	return buf.String()
}

// AsExitStatus looks at an Execute() error and decides
// how the process should finish
func AsExitStatus(err error) (int, error) {
	switch err {
	case nil:
		// all good
		return ExitStatusOK, nil
	case pflag.ErrHelp:
		// not an error, just an exit condition
		return ExitStatusMinor, nil
	default:
		// ExitStatus wrapper
		if e, ok := err.(*ErrorExitCode); ok {
			return e.Code, e.Err
		}

		if e, ok := err.(interface {
			ExitStatus() int
		}); ok {
			// Error knows the right ExitStatus, but it
			// could include more information
			return e.ExitStatus(), err
		}

		// Some other error that doesn't include ExitStatus
		return ExitStatusMajor, err
	}
}
