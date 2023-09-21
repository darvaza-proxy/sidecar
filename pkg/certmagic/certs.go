package certmagic

import (
	"crypto/tls"
	"crypto/x509"

	"darvaza.org/core"
	"darvaza.org/darvaza/shared/storage"
)

var (
	_ storage.Store = (*Store)(nil)
)

// GetCAPool provides RootCAs and ClientCSs for tls.Config
func (*Store) GetCAPool() *x509.CertPool {
	panic(core.ErrNotImplemented)
}

// GetCertificate implements tls.Config.GetCertificate
func (*Store) GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return nil, core.ErrNotImplemented
}

// AddTrustedRoot ...
func (*Store) AddTrustedRoot(_ string) (bool, error) {
	return false, core.ErrNotImplemented
}

// SetIssuerKey ...
func (*Store) SetIssuerKey(_ string, _ string) error {
	return core.ErrNotImplemented
}

// SetKey ...
func (*Store) SetKey(_ string, _ ...string) error {
	return core.ErrNotImplemented
}
