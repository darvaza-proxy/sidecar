package reconnect

import (
	"context"
	"net"
	"time"

	"darvaza.org/core"
	"darvaza.org/slog"
	"darvaza.org/slog/handlers/discard"
	"darvaza.org/x/config"
)

// OptionFunc is a function used on [New] to modify
// the [Options] before the [Client] is created
type OptionFunc func(*Options) error

// Options describe how the [Client] will operate
type Options struct {
	Context context.Context
	Logger  slog.Logger

	KeepAlive      time.Duration `default:"5s"`
	ConnectTimeout time.Duration `default:"2s"`
	ReadTimeout    time.Duration `default:"2s"`
	WriteTimeout   time.Duration `default:"2s"`

	ReconnectDelay time.Duration `default:"5s"`
	WaitReconnect  Waiter

	// OnSession is expected to block until it's done.
	OnSession func(context.Context) error
	// OnDisconnect is called after closing the connection and can be used to
	// prevent further connection retries.
	OnDisconnect func(context.Context, net.Conn) error
	// OnError is called after all errors and gives us the opportunity to
	// decide how the error should be treated by the reconnection logic.
	OnError func(context.Context, net.Conn, error) error
}

// SetDefaults fills gaps in [Options]
func (opts *Options) SetDefaults() error {
	if err := config.Set(opts); err != nil {
		return err
	}

	if opts.Context == nil {
		opts.Context = context.Background()
	}

	if opts.WaitReconnect == nil {
		opts.SetConstantWaitReconnect(opts.ReconnectDelay)
	}

	if opts.Logger == nil {
		opts.Logger = discard.New()
	}

	return nil
}

// SetConstantWaitReconnect sets [Options] reconnection rate.
func (opts *Options) SetConstantWaitReconnect(d time.Duration) *Options {
	opts.WaitReconnect = NewConstantWaiter(d)
	opts.ReconnectDelay = core.IIf(d < 0, 0, d)

	return opts
}

// New creates a Client using the Options
func (opts *Options) New(network, address string, funcs ...OptionFunc) (*Client, error) {
	// set defaults
	if err := opts.SetDefaults(); err != nil {
		return nil, err
	}

	// apply alterations
	for _, fn := range funcs {
		if err := fn(opts); err != nil {
			return nil, err
		}
	}

	return New(network, address, opts)
}
