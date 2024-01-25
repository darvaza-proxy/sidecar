package service

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/spf13/pflag"
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
		var code int
		var wrapped error

		switch e := err.(type) {
		case interface{ ExitStatus() int }:
			code = e.ExitStatus()
			wrapped = errors.Unwrap(err)
		case interface{ ExitCode() int }:
			code = e.ExitCode()
			wrapped = errors.Unwrap(err)
		default:
			code = ExitStatusMajor
		}

		if wrapped != nil {
			err = wrapped
		}

		return code, err
	}
}
