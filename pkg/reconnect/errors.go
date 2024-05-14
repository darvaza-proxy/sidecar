package reconnect

import (
	"io/fs"
	"syscall"

	"darvaza.org/core"

	"darvaza.org/sidecar/pkg/utils"
)

var (
	// ErrAbnormalConnect indicates the dialer didn't return error
	// nor connection.
	ErrAbnormalConnect = core.Wrap(syscall.ECONNABORTED, "abnormal response")

	// ErrNotConnected indicates the [Client] isn't currently connected.
	ErrNotConnected = core.Wrap(fs.ErrClosed, "connection closed")
)

// IsFatal tells if the error means the connection
// should be closed and not retried.
func IsFatal(err error) bool {
	is, _ := core.IsErrorFn2(checkIsFatal, err)
	return is
}

func checkIsFatal(err error) (is bool, ok bool) {
	switch err {
	case ErrDoNotReconnect:
		// do-not reconnect
		return true, true
	default:
		// noise errors are never fatal
		if is, _ := checkNoiseError(err); is {
			return false, true
		}

		// reconnect if temporary
		return utils.CheckIsTemporary(err)
	}
}

// IsNoiseError tells if the error can be ignored.
func IsNoiseError(err error) bool {
	is, _ := core.IsErrorFn2(checkNoiseError, err)
	return is
}

func checkNoiseError(err error) (is, ok bool) {
	switch err {
	case nil,
		fs.ErrClosed,
		syscall.ECONNABORTED,
		syscall.ECONNREFUSED,
		syscall.ECONNRESET:
		// reconnect
		return true, true
	default:
		// unknown
		return false, false
	}
}
