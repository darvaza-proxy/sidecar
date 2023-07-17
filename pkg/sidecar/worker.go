package sidecar

import (
	"context"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

// onWorkerError is called by srv.wg on errors
func (srv *Server) onWorkerError(err error) error {
	// TODO: decide when is an error worth calling
	// srv.Fail(err)
	srv.error(err).Print("onWorkerError")

	return nil
}

// Shutdown initiates a shutdown of all workers with optional
// fatal timeout while waiting for the workers to finish.
func (srv *Server) Shutdown(timeout time.Duration) error {
	var ok atomic.Bool

	// once srv.Wait() finishes, we are done
	defer ok.Store(true)

	srv.tryCancel(nil)

	if timeout > 0 {
		time.AfterFunc(timeout, func() {
			if !ok.Load() {
				srv.fatal(nil).Print("graceful shutdown timed out")
			}
		})
	}

	return srv.Wait()
}

// Cancel initiates a shutdown of all workers
func (srv *Server) Cancel() {
	srv.tryCancel(nil)
}

// Fail initiates a shutdown with a reason
func (srv *Server) Fail(err error) {
	srv.tryCancel(err)
}

func (srv *Server) tryCancel(err error) {
	// once
	if srv.cancelled.CompareAndSwap(false, true) {
		if err != nil {
			srv.err.Store(err)
		}
		srv.cancel()
	}
}

// Cancelled tells if the server has been cancelled
func (srv *Server) Cancelled() bool {
	return srv.cancelled.Load()
}

// Err returns the reasons of the shutdown, if any
func (srv *Server) Err() error {
	if err, ok := srv.err.Load().(error); ok {
		return err
	} else if srv.Cancelled() {
		return os.ErrClosed
	} else {
		return nil
	}
}

// Wait blocks until all workers have exited
func (srv *Server) Wait() error {
	srv.wg.Wait()

	err := srv.Err()
	switch err {
	case nil, os.ErrClosed:
		// no error
		return nil
	default:
		// actual error
		return err
	}
}

// Spawn starts the initial workers
func (srv *Server) Spawn(h http.Handler, healthy time.Duration) error {
	var ok bool

	defer func() {
		if !ok {
			srv.Cancel()
		}
	}()

	if err := srv.spawnHTTPServer(h); err != nil {
		return err
	}

	if healthy > 0 {
		time.Sleep(healthy)
	}

	ok = true
	return srv.Err()
}

// Go runs a worker on the Server's Context
func (srv *Server) Go(fn func(ctx context.Context) error) {
	srv.wg.Go(func() error {
		return fn(srv.ctx)
	})
}
