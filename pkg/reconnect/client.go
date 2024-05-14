// Package reconnect implement a generic retrying TCP client
package reconnect

import (
	"bufio"
	"context"
	"net"
	"sync"
	"syscall"
	"time"

	"darvaza.org/slog"
)

// Client is a reconnecting TCP Client.
type Client struct {
	mu sync.Mutex

	ctx     context.Context
	cancel  context.CancelCauseFunc
	running bool
	err     error
	wg      sync.WaitGroup

	readTimeout  time.Duration
	writeTimeout time.Duration

	onSession     func(context.Context) error
	onDisconnect  func(context.Context, net.Conn) error
	onError       func(context.Context, net.Conn, error) error
	waitReconnect func(context.Context) error

	logger  slog.Logger
	dialer  net.Dialer
	network string
	address string

	conn net.Conn
	in   *bufio.Reader
	out  *bufio.Writer
}

// New creates a new [Client] using the provided [Options]
func New(network, address string, opts *Options) (*Client, error) {
	if opts == nil {
		opts = new(Options)
	}

	if err := opts.SetDefaults(); err != nil {
		return nil, err
	}

	c := &Client{
		dialer: net.Dialer{
			Timeout:   opts.ConnectTimeout,
			KeepAlive: opts.KeepAlive,
		},

		logger:        opts.Logger,
		readTimeout:   opts.ReadTimeout,
		writeTimeout:  opts.WriteTimeout,
		onSession:     opts.OnSession,
		onDisconnect:  opts.OnDisconnect,
		onError:       opts.OnError,
		waitReconnect: opts.WaitReconnect,

		network: network,
		address: address,
	}

	c.ctx, c.cancel = context.WithCancelCause(opts.Context)
	return c, nil
}

// Connect launches the [Client]
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return syscall.EBUSY
	}

	conn, err := c.dial()
	if err = c.handlePossiblyFatalError(conn, err, ""); err != nil {
		return err
	}

	c.wg.Add(1)
	c.unsafeSetConn(conn)
	c.running = true
	c.err = nil

	go func() {
		defer c.wg.Done()

		c.run()
	}()

	return nil
}
