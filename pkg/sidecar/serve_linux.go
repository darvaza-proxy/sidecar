//go:build linux

package sidecar

import (
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudflare/tableflip"
)

// ListenAndServe runs the Server
func (srv *Server) ListenAndServe(app http.Handler) error {
	// Upgrader
	upg, err := tableflip.New(tableflip.Options{
		PIDFile: srv.cfg.Supervision.PIDFile,
	})
	if err != nil {
		return err
	}

	// Prepare Server
	if app == nil {
		app = http.NotFoundHandler()
	}

	// Listen
	if err := srv.ListenWithUpgrader(upg); err != nil {
		return err
	}

	// Watch signals
	go srv.watchSignals(upg, app)

	// Close app before exiting
	if r, ok := app.(io.Closer); ok {
		defer r.Close()
	}

	// Attempt to start the server
	if err := srv.Spawn(app, srv.cfg.Supervision.HealthWait); err != nil {
		return err
	}

	// Notify being ready for service
	if err := upg.Ready(); err != nil {
		return err
	}
	<-upg.Exit()

	// Wait for connections to drain.
	return srv.Shutdown(srv.cfg.Supervision.GracefulTimeout)
}

func (srv *Server) watchSignals(upg *tableflip.Upgrader, app http.Handler) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)

	for signum := range sig {
		switch signum {
		case syscall.SIGHUP:
			srv.onSIGHUP(upg, app)
		case syscall.SIGUSR2:
			srv.onSIGUSR2(upg, app)
		case syscall.SIGINT, syscall.SIGTERM:
			srv.onSIGTERM(upg, app)
		}
	}
}

func (srv *Server) onSIGHUP(upg *tableflip.Upgrader, app http.Handler) {
	// attempt to reload config on SIGHUP if supported
	// or an upgrade if it isn't
	if r, ok := app.(Reloader); ok {
		if err := r.Reload(); err != nil {
			srv.error(err).Println("reload failed")
		}
	} else if err := upg.Upgrade(); err != nil {
		srv.error(err).Println("upgrade failed")
	}
}

func (srv *Server) onSIGUSR2(upg *tableflip.Upgrader, _ http.Handler) {
	// attempt to upgrade on SIGUSR2
	if err := upg.Upgrade(); err != nil {
		srv.error(err).Println("upgrade failed")
	}
}

func (srv *Server) onSIGTERM(upg *tableflip.Upgrader, _ http.Handler) {
	// terminate on SIGINT or SIGTERM
	srv.warn().Println("terminate signal received")
	upg.Stop()
}
