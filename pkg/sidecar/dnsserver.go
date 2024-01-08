package sidecar

import (
	dns "darvaza.org/sidecar/pkg/sidecar/dnsserver"
)

func (srv *Server) initDNSServer() error {
	dsc, ok := srv.newDNSServerConfig()
	if ok {
		ds, err := dsc.New(&srv.eg)
		if err != nil {
			return err
		}
		srv.ds = ds
	}

	return nil
}

func (srv *Server) newDNSServerConfig() (*dns.Config, bool) {
	dc := &srv.cfg.DNS
	if !dc.Enabled {
		return nil, false
	}

	mTLS := srv.cfg.DNS.MutualTLSOnly

	dsc := &dns.Config{
		Logger:  srv.cfg.Logger,
		Context: srv.cfg.Context,

		// TLS
		TLSConfig: srv.newTLSServerConfig(mTLS),

		// Address
		Bind: dns.BindingConfig{
			Addrs: srv.cfg.Addresses.Addresses,

			// Ports
			Port:    dc.Port,
			TLSPort: dc.TLSPort,
		},

		// DNS
		MaxTCPQueries: dc.MaxTCPQueries,
		ReadTimeout:   dc.ReadTimeout,
		IdleTimeout:   dc.IdleTimeout,

		GracefulTimeout: srv.cfg.Supervision.GracefulTimeout,
	}

	return dsc, true
}
