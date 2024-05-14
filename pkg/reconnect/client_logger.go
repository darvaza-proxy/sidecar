package reconnect

import (
	"fmt"
	"net"

	"darvaza.org/core"
	"darvaza.org/slog"

	"darvaza.org/sidecar/pkg/utils"
)

// WithInfo ...
func (c *Client) WithInfo(addr net.Addr) (slog.Logger, bool) {
	l, ok := c.logger.Info().WithEnabled()
	if !ok {
		return nil, false
	}

	l = utils.LogWithAddress(l, "", addr)

	return l, true
}

// WithError ...
func (c *Client) WithError(addr net.Addr, err error) (slog.Logger, bool) {
	l, ok := c.logger.Error().WithEnabled()
	if !ok {
		return nil, false
	}

	l = utils.LogWithAddress(l, "", addr)
	l = utils.LogWithError(l, err)

	return l, true
}

// SayRemote ...
func (c *Client) SayRemote(conn net.Conn, note string, args ...any) {
	var ra net.Addr
	if conn != nil {
		ra = conn.RemoteAddr()
	}

	if l, ok := c.WithInfo(ra); ok {
		if len(args) > 0 {
			l.Printf(note, args...)
		} else {
			l.Print(note)
		}
	}
}

// SayRemoteError ...
func (c *Client) SayRemoteError(conn net.Conn, err error, note string, args ...any) {
	var ra net.Addr
	if conn != nil {
		ra = conn.RemoteAddr()
	}

	if l, ok := c.WithError(ra, err); ok {
		if len(args) > 0 {
			l.Printf(note, args...)
		} else {
			l.Print(note)
		}
	}
}

// SayError ...
func (c *Client) SayError(err error) {
	c.SayRemoteError(nil, err, "")
}

// Panic ...
func (c *Client) Panic(err error, note string, args ...any) {
	defer func() { _ = c.terminate(err) }()

	if err != nil {
		err = core.Wrap(err, note, args...)
		note = ""
	} else if len(args) > 0 {
		note = fmt.Sprintf(note, args...)
	}

	l := c.logger.Panic()
	l = utils.LogWithError(l, err)
	l.Panic().Print(note)

	panic("unreachable")
}
