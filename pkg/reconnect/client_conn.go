package reconnect

import (
	"bufio"
	"net"
	"time"

	"darvaza.org/core"
	"darvaza.org/sidecar/pkg/utils"
)

var (
	_ Reader  = (*Client)(nil)
	_ Writer  = (*Client)(nil)
	_ Flusher = (*Client)(nil)
	_ Closer  = (*Client)(nil)
)

// ResetDeadline sets the connection's read and write deadlines using
// the default values.
func (c *Client) ResetDeadline() error {
	return c.SetDeadline(c.readTimeout, c.writeTimeout)
}

// ResetReadDeadline resets the connection's read deadline using
// the default duration.
func (c *Client) ResetReadDeadline() error {
	t := utils.TimeoutToAbsoluteTime(time.Now(), c.readTimeout)

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return ErrNotConnected
	}

	return c.conn.SetReadDeadline(t)
}

// ResetWriteDeadline resets the connection's write deadline using
// the default duration.
func (c *Client) ResetWriteDeadline() error {
	t := utils.TimeoutToAbsoluteTime(time.Now(), c.writeTimeout)

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return ErrNotConnected
	}

	return c.conn.SetWriteDeadline(t)
}

// SetDeadline sets the connections's read and write deadlines.
// if write is zero but read is positive, write is set using the same
// value as read.
// zero or negative can be used to disable the deadline.
func (c *Client) SetDeadline(read, write time.Duration) error {
	if read > 0 && write == 0 {
		write = read
	}

	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return ErrNotConnected
	}

	t := utils.TimeoutToAbsoluteTime(now, read)
	if err := c.conn.SetReadDeadline(t); err != nil {
		return err
	}

	t = utils.TimeoutToAbsoluteTime(now, write)
	return c.conn.SetWriteDeadline(t)
}

// dials attempts to stablish a connection to the server.
func (c *Client) dial() (net.Conn, error) {
	conn, err := c.dialer.DialContext(c.ctx, c.network, c.address)
	if conn == nil {
		if err == nil {
			err = &net.OpError{
				Op:  "dial",
				Net: c.network,
				Err: core.Wrap(ErrAbnormalConnect, c.address),
			}
		}
		return nil, err
	}

	return conn, nil
}

// reconnect waits before dialing.
func (c *Client) reconnect() (net.Conn, error) {
	if fn := c.waitReconnect; fn != nil {
		if err := fn(c.ctx); err != nil {
			return nil, err
		}
	}

	return c.dial()
}

// Read implements a buffered io.Reader
func (c *Client) Read(p []byte) (int, error) {
	c.mu.Lock()
	r := c.in
	c.mu.Unlock()

	if r == nil {
		return 0, ErrNotConnected
	}

	return r.Read(p)
}

// Write implements a buffered io.Writer
// warrantied to buffer all the given data or fail.
func (c *Client) Write(p []byte) (int, error) {
	c.mu.Lock()
	w := c.out
	c.mu.Unlock()

	if w == nil {
		return 0, ErrNotConnected
	}

	total := 0
	for len(p) > 0 {
		n, err := w.Write(p)

		switch {
		case err != nil:
			return total, err
		default:
			total += n
			p = p[n:]
		}
	}

	return total, nil
}

// Flush blocks until all the buffered output
// has been written, or an error occurs.
func (c *Client) Flush() error {
	c.mu.Lock()
	w := c.out
	c.mu.Unlock()

	if w == nil {
		return ErrNotConnected
	}

	return w.Flush()
}

func (c *Client) setConn(conn net.Conn) net.Conn {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.unsafeSetConn(conn)
}

func (c *Client) unsafeSetConn(conn net.Conn) (prev net.Conn) {
	prev, c.conn = c.conn, conn
	if conn == nil {
		c.in, c.out = nil, nil
	} else {
		c.in = bufio.NewReader(conn)
		c.out = bufio.NewWriter(conn)
	}

	return prev
}
