//go:build windows

package sidecar

import (
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// ListenAndServe runs the Server
func (srv *Server) ListenAndServe(app http.Handler) error {
	// Prepare Server
	if app == nil {
		app = http.NotFoundHandler()
	}

	// Listen
	if err := srv.Listen(); err != nil {
		return err
	}

	// Watch signals
	go srv.watchSignals(app)

	// Close app before exiting
	if r, ok := app.(io.Closer); ok {
		defer r.Close()
	}

	// Attempt to start the server
	if err := srv.Spawn(app, srv.cfg.Supervision.HealthWait); err != nil {
		return err
	}
	<-srv.ctx.Done()

	// Wait for connections to drain.
	return srv.Shutdown(srv.cfg.Supervision.GracefulTimeout)
}

func (srv *Server) watchSignals(app http.Handler) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	for signum := range sig {
		switch signum {
		case syscall.SIGHUP:
			srv.onSIGHUP(app)
		case syscall.SIGINT, syscall.SIGTERM:
			srv.onSIGTERM(app)
		}
	}
}

func (srv *Server) onSIGHUP(app http.Handler) {
	// attempt to reload config on SIGHUP if supported
	if r, ok := app.(Reloader); ok {
		if err := r.Reload(); err != nil {
			srv.error(err).Println("reload failed")
		}
	}
}

func (srv *Server) onSIGTERM(_ http.Handler) {
	// terminate on SIGINT or SIGTERM
	srv.warn().Println("terminate signal received")
	srv.cancel()
}
