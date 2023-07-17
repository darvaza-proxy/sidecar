package sidecar

import (
	"crypto/tls"
	"crypto/x509"

	"darvaza.org/darvaza/shared/storage"
	"darvaza.org/darvaza/shared/storage/simple"
)

func newTLSStore(cfg *Config) (storage.Store, error) {
	sc := &simple.Config{
		Logger: cfg.Logger,
	}

	return sc.New(cfg.TLS.Key,
		cfg.TLS.Cert,
		cfg.TLS.Roots)
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
