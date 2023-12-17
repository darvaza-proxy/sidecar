package sidecar

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"
)

// onWorkerError is called by srv.eg on errors
func (srv *Server) onWorkerError(err error) {
	srv.error(err).Print("onWorkerError")
}

// Shutdown initiates a shutdown of all workers with optional
// fatal timeout while waiting for the workers to finish.
func (srv *Server) Shutdown(timeout time.Duration) error {
	srv.eg.Cancel(context.Canceled)

	if timeout > 0 {
		select {
		case <-time.After(timeout):
			// timed out
			return errors.New("graceful shutdown timed out")
		case <-srv.eg.Done():
			// finished
		}
	}

	return srv.Wait()
}

// Cancel initiates a shutdown of all workers
func (srv *Server) Cancel() {
	srv.eg.Cancel(nil)
}

// Fail initiates a shutdown with a reason
func (srv *Server) Fail(cause error) {
	srv.eg.Cancel(cause)
}

// IsCancelled tells if the server has been cancelled
func (srv *Server) IsCancelled() bool {
	return srv.eg.IsCancelled()
}

// Err returns the reasons of the shutdown, if any
func (srv *Server) Err() error {
	if err := srv.eg.Err(); err != nil {
		return err
	}

	if srv.IsCancelled() {
		return os.ErrClosed
	}

	return nil
}

// Wait blocks until all workers have exited
func (srv *Server) Wait() error {
	err := srv.eg.Wait()

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
	srv.spawnHTTPServer(h)

	if healthy > 0 {
		select {
		case <-time.After(healthy):
			// done waiting
		case <-srv.eg.Cancelled():
			// failed while waiting, let them flush out.
			return srv.eg.Wait()
		}
	}

	return nil
}

// Go runs a worker on the Server's Context
func (srv *Server) Go(run func(ctx context.Context) error) {
	srv.eg.Go(run, nil)
}

// GoWithShutdown runs a worker on the Server's Context, and a shutdown sentinel.
func (srv *Server) GoWithShutdown(run func(context.Context) error, shutdown func() error) {
	srv.eg.Go(run, shutdown)
}
