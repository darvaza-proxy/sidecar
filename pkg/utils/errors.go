package utils

import (
	"darvaza.org/core"
)

// IsTimeout recursively tests if the error represents a Time-out condition.
func IsTimeout(err error) bool {
	is, _ := core.IsErrorFn2(CheckIsTimeout, err)
	return is
}

// CheckIsTimeout non-recursively if the error is a
// time-out condition without unwrapping.
func CheckIsTimeout(err error) (is, ok bool) {
	switch e := err.(type) {
	case nil:
		return false, true
	case interface {
		IsTimeout() bool
	}:
		return e.IsTimeout(), true
	case interface {
		Timeout() bool
	}:
		return e.Timeout(), true
	default:
		// unknown
		return false, false
	}
}

// IsTemporary recursively tests if the error is temporary
func IsTemporary(err error) bool {
	is, _ := core.IsErrorFn2(CheckIsTemporary, err)
	return is
}

// CheckIsTemporary tests if the error is temporary
// without unwrapping.
func CheckIsTemporary(err error) (is, ok bool) {
	switch e := err.(type) {
	case nil:
		return false, true
	case interface {
		IsTemporary() bool
	}:
		return e.IsTemporary(), true
	case interface {
		Temporary() bool
	}:
		return e.Temporary(), true
	default:
		// time-out conditions are temporary
		return CheckIsTimeout(err)
	}
}
