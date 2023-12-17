package sidecar

import (
	"crypto/tls"
	"crypto/x509"

	"darvaza.org/core"
	"darvaza.org/darvaza/shared/storage"
)

func newFallbackTLSStore() (storage.Store, error) {
	// TODO: self-signed
	return nil, core.ErrNotImplemented
}

// revive:disable:flag-parameter
func (srv *Server) newTLSServerConfig(mTLS bool) *tls.Config {
	// revive:enable:flag-parameter
	var rootCAs, clientCAs *x509.CertPool

	rootCAs = srv.tls.GetCAPool()
	if mTLS {
		clientCAs = rootCAs
	}

	return &tls.Config{
		ServerName: srv.cfg.Name,

		GetCertificate: srv.tls.GetCertificate,
		ClientCAs:      clientCAs,
		RootCAs:        rootCAs,
	}
}
