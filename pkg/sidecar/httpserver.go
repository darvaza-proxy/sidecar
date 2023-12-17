package sidecar

import (
	"darvaza.org/sidecar/pkg/sidecar/httpserver"
)

func (srv *Server) initHTTPServer() error {
	hsc := srv.newHTTPServerConfig()
	hs, err := hsc.New(&srv.eg)
	if err != nil {
		return err
	}
	srv.hs = hs
	return nil
}

func (srv *Server) newHTTPServerConfig() *httpserver.Config {
	mTLS := srv.cfg.HTTP.MutualTLSOnly

	hsc := &httpserver.Config{
		Logger:  srv.cfg.Logger,
		Context: srv.cfg.Context,

		// TLS
		TLSConfig: srv.newTLSServerConfig(mTLS),

		// Addresses
		Bind: httpserver.BindingConfig{
			Addrs: srv.cfg.Addresses.Addresses,

			Port:          srv.cfg.HTTP.Port,
			PortInsecure:  srv.cfg.HTTP.PortInsecure,
			AllowInsecure: srv.cfg.HTTP.EnableInsecure,
		},

		// HTTP
		ReadTimeout:       srv.cfg.HTTP.ReadTimeout,
		ReadHeaderTimeout: srv.cfg.HTTP.ReadHeaderTimeout,
		WriteTimeout:      srv.cfg.HTTP.WriteTimeout,
		IdleTimeout:       srv.cfg.HTTP.IdleTimeout,

		GracefulTimeout: srv.cfg.Supervision.GracefulTimeout,
	}

	return hsc
}
