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

func (srv *Server) getGetCertificateForServer() func(*tls.ClientHelloInfo) (*tls.Certificate,
	error) {
	return srv.tls.GetCertificate
}

func (*Server) getRootCAsForServer() func() *x509.CertPool {
	return nil
}

func (srv *Server) getClientCAsForServer() func() *x509.CertPool {
	if srv.cfg.HTTP.MutualTLSOnly {
		return srv.tls.GetCAPool
	}
	return nil
}
