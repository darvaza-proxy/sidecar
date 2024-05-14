package reconnect

import (
	"context"
	"net"
	"syscall"

	"darvaza.org/core"
)

// Wait blocks until the [Client] workers have finished,
// and returns the cancellation reason.
func (c *Client) Wait() error {
	c.wg.Wait()
	return c.Err()
}

// Done returns a channel that watches the [Client] workers,
// and provides the cancellation reason.
func (c *Client) Done() <-chan error {
	var barrier chan error
	go func() {
		defer close(barrier)
		barrier <- c.Wait()
	}()

	return barrier
}

// Err returns the cancellation reason
func (c *Client) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return filterError(c.err)
}

// filterError removes non-error errors from Err()
func filterError(err error) error {
	if err != nil {
		if is, _ := core.IsErrorFn2(isNonError, err); is {
			return nil
		}
	}
	return err
}

func isNonError(err error) (is, ok bool) {
	switch err {
	case context.Canceled, ErrDoNotReconnect:
		return true, true
	default:
		return false, false
	}
}

// Close initiates a shutdown
func (c *Client) Close() error {
	return c.terminate(nil)
}

// Shutdown initiates a shutdown and wait until the workers
// are done, or the given context times out.
func (c *Client) Shutdown(ctx context.Context) error {
	_ = c.terminate(nil)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.Done():
		return c.Err()
	}
}

func (c *Client) run() {
	var conn net.Conn
	var abort bool

	defer func() {
		if err := core.AsRecovered(recover()); err != nil {
			_ = c.doOnError(nil, err, "reconnect.Client")
			_ = c.terminate(err)
		}

		c.mu.Lock()
		c.running = false
		c.mu.Unlock()
	}()

	for !abort {
		conn, abort = c.doConnect()

		if conn != nil {
			e1 := c.runConn(conn)
			e2 := c.doDisconnect()

			abort = c.runError(conn, e1, e2)
		}
	}
}

func (c *Client) doConnect() (net.Conn, bool) {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()

	if conn != nil {
		return conn, false // ready
	}

	conn, err := c.reconnect()
	if abort := c.runError(nil, err, nil); abort {
		return nil, true // abort
	}

	if conn != nil {
		c.setConn(conn)
	}

	return conn, false // ready
}

func (c *Client) doDisconnect() error {
	conn := c.setConn(nil)

	c.SayRemote(conn, "disconnected")
	if c.onDisconnect != nil {
		return c.onDisconnect(c.ctx, conn)
	}
	return nil
}

func (c *Client) runConn(conn net.Conn) error {
	defer conn.Close()

	c.SayRemote(conn, "connected")
	if err := c.ResetDeadline(); err != nil {
		return err
	}

	if c.onSession != nil {
		var catch core.Catcher
		return catch.Try(func() error {
			return c.onSession(c.ctx)
		})
	}

	return nil
}

func (c *Client) runError(conn net.Conn, e1, e2 error) bool {
	e1 = c.handlePossiblyFatalError(conn, e1, "")
	e2 = c.handlePossiblyFatalError(conn, e2, "")
	return e1 != nil || e2 != nil
}

func (c *Client) terminate(cause error) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.err == nil {
		if cause == nil {
			cause = context.Canceled
		}

		c.err = cause
		c.cancel(cause)
	}

	return filterError(c.err)
}

// Go spawns a goroutine within the [Client]'s context.
func (c *Client) Go(name string, fn func(context.Context) error) {
	if fn == nil {
		c.Panic(syscall.EINVAL, "Client.Go called without a handler")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.wg.Add(1)
	go func() {
		var catcher core.Catcher

		defer c.wg.Done()

		err := catcher.Do(func() error {
			return fn(c.ctx)
		})

		if c.handlePossiblyFatalError(nil, err, name) != nil {
			return
		}
	}()
}

// handlePossiblyFatalError handles an error and returns nil if it wasn't fatal.
// fatal errors should terminate the worker immediately.
// the returned error is unfiltered.
func (c *Client) handlePossiblyFatalError(conn net.Conn, err error, note string) error {
	if err != nil {
		err = c.doOnError(conn, err, note)
		if IsFatal(err) {
			_ = c.terminate(err)
			return err // unfiltered
		}
	}
	return nil
}

func (c *Client) doOnError(remoteConn net.Conn, err error, note string, args ...any) error {
	if err != nil {
		c.SayRemoteError(remoteConn, err, note, args...)
		if c.onError != nil {
			err = c.onError(c.ctx, remoteConn, err)
		}
	}
	return err
}
