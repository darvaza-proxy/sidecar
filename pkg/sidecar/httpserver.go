package sidecar

import (
	"net/http"

	"darvaza.org/darvaza/agent/httpserver"
)

func (srv *Server) newHTTPServerConfig() *httpserver.Config {
	hsc := &httpserver.Config{
		Logger:  srv.cfg.Logger,
		Context: srv.ctx,

		// Addresses
		Bind: httpserver.BindingConfig{
			Interfaces: srv.cfg.Addresses.Interfaces,
			Addresses:  srv.cfg.Addresses.Addresses,

			PortInsecure:  srv.cfg.HTTP.Port,
			AllowInsecure: true,
		},

		// HTTP
		ReadTimeout:       srv.cfg.HTTP.ReadTimeout,
		ReadHeaderTimeout: srv.cfg.HTTP.ReadHeaderTimeout,
		WriteTimeout:      srv.cfg.HTTP.WriteTimeout,
		IdleTimeout:       srv.cfg.HTTP.IdleTimeout,
	}

	return hsc
}

func (srv *Server) spawnHTTPServer(h http.Handler) error {
	srv.wg.Go(func() error {
		return srv.hs.Serve(h)
	})

	srv.wg.Go(func() error {
		<-srv.ctx.Done()
		srv.hs.Cancel()
		return srv.hs.Wait()
	})

	return nil
}
