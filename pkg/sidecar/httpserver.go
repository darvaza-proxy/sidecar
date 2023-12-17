package sidecar

import (
	"context"
	"net/http"

	"darvaza.org/darvaza/agent/httpserver"
)

func (srv *Server) initHTTPServer() error {
	hsc := srv.newHTTPServerConfig()
	hs, err := hsc.New()
	if err != nil {
		return err
	}
	srv.hs = hs
	return nil
}

func (srv *Server) newHTTPServerConfig() *httpserver.Config {
	da := &srv.cfg.Addresses
	addrs := make([]string, 0, len(da.Addresses))
	for _, addr := range da.Addresses {
		addrs = append(addrs, addr.String())
	}

	hsc := &httpserver.Config{
		Logger:  srv.cfg.Logger,
		Context: srv.cfg.Context,

		// Addresses
		Bind: httpserver.BindingConfig{
			Addresses: addrs,

			Port:          srv.cfg.HTTP.Port,
			PortInsecure:  srv.cfg.HTTP.PortInsecure,
			AllowInsecure: srv.cfg.HTTP.EnableInsecure,
		},

		// HTTP
		ReadTimeout:       srv.cfg.HTTP.ReadTimeout,
		ReadHeaderTimeout: srv.cfg.HTTP.ReadHeaderTimeout,
		WriteTimeout:      srv.cfg.HTTP.WriteTimeout,
		IdleTimeout:       srv.cfg.HTTP.IdleTimeout,

		// TLS
		GetCertificate: srv.getGetCertificateForServer(),
		GetRootCAs:     srv.getRootCAsForServer(),
		GetClientCAs:   srv.getClientCAsForServer(),
	}

	return hsc
}

func (srv *Server) spawnHTTPServer(h http.Handler) {
	srv.eg.Go(func(_ context.Context) error {
		return srv.hs.Serve(h)
	}, func() error {
		srv.hs.Cancel()
		return srv.hs.Wait()
	})
}
