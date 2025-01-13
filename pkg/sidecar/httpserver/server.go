// Package httpserver implements a HTTP/HTTPS/HTTP3 server
// for darvaza sidecars.
package httpserver

import (
	"context"
	"net/http"
	"sync"
	"syscall"
	"time"

	"darvaza.org/core"
)

// Server is an HTTP/1, HTTP/2, HTTP/3 server built
// around a shared [core.ErrGroup].
type Server struct {
	mu  sync.Mutex
	cfg Config

	eg *core.ErrGroup
	sl *Listeners

	quicAltSvc string
}

// Spawn starts all workers and optionally waits a given amount
// to make sure they didn't fail.
func (srv *Server) Spawn(h http.Handler, wait time.Duration) error {
	if h == nil {
		h = http.NotFoundHandler()
	}

	graceful := srv.cfg.GracefulTimeout

	for _, fn := range []func() error{
		func() error { return srv.prepare() },
		func() error { return srv.spawnH2C(h, srv.sl.Insecure, graceful) },
		func() error { return srv.spawnH2(h, srv.sl.Secure, graceful) },
		func() error { return srv.spawnH3(h, srv.sl.QUIC, graceful) },
	} {
		if err := fn(); err != nil {
			srv.eg.Cancel(err)
			_ = srv.sl.Close()
			return srv.eg.Err()
		}
	}

	if wait > 0 {
		select {
		case <-time.After(wait):
			// done waiting
		case <-srv.eg.Cancelled():
			// failed while waiting.
			return srv.eg.Wait()
		}
	}

	return srv.eg.Err()
}

func (srv *Server) prepare() error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	switch {
	case srv.sl == nil:
		// not listening
		return core.Wrap(core.ErrInvalid, "no listeners available")
	case srv.sl.up.CompareAndSwap(false, true):
		// once
		if srv.eg == nil {
			srv.eg = &core.ErrGroup{
				Parent: srv.cfg.Context,
			}
		}
		return nil
	default:
		// not again
		return core.Wrap(syscall.EBUSY, "server already running")
	}
}

// Serve starts all workers and waits until they have
// finished.
func (srv *Server) Serve(h http.Handler) error {
	if err := srv.Spawn(h, 0); err != nil {
		// failed to prepare
		return err
	}

	return srv.Wait()
}

// Cancel initiates a cancellation with the given
// reason.
func (srv *Server) Cancel(cause error) {
	srv.eg.Cancel(cause)
}

// Wait waits until all workers have finished.
func (srv *Server) Wait() error {
	return srv.eg.Wait()
}

// Close tries to initiate a cancellation
// if the server is running, and closes all listeners.
// Errors are ignored.
func (srv *Server) Close() error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.sl == nil {
		// nothing to close
		return nil
	}

	if srv.sl.up.Load() {
		// running, cancel first.
		srv.eg.Cancel(context.Canceled)
	}

	// Close listeners
	return srv.sl.Close()
}
