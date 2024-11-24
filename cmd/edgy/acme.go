package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"darvaza.org/core"
	"darvaza.org/slog"
	"darvaza.org/x/tls/store/autocert"
)

// ACMEConfig ...
type ACMEConfig struct {
	URL   string
	Token string
	EMail string

	Roots string // TODO: use
	Cert  string // TODO: use
	Key   string // TODO: use

	CacheDir string // TODO: use
}

// SetDefaults ...
func (ac *ACMEConfig) SetDefaults() error {
	if ac.URL == "" {
		ac.URL = autocert.LetsEncryptStagingURL
	}
	return nil
}

// Export ...
func (ac *ACMEConfig) Export(logger slog.Logger) (*autocert.Config, error) {
	if ac == nil {
		return nil, core.ErrNilReceiver
	}

	cfg := &autocert.Config{
		Logger:   logger,
		CacheDir: ac.CacheDir,

		DirectoryURL: ac.URL,
		AcceptTOS:    true,
		EMail:        ac.EMail,
		BearerToken:  ac.Token,

		TrustedCAs: ac.exportCACertPool(),
		ClientCert: ac.exportClientCert(),
	}

	return cfg, nil
}

func (ac *ACMEConfig) exportCACertPool() *x509.CertPool
func (ac *ACMEConfig) exportClientCert() *tls.Certificate

// New ...
func (ac *ACMEConfig) New(ctx context.Context, logger slog.Logger) (*autocert.Store, error) {
	if ac == nil {
		return nil, core.ErrNilReceiver
	}

	err := ac.SetDefaults()
	if err != nil {
		return nil, err
	}

	cfg, err := ac.Export(logger)
	if err != nil {
		return nil, err
	}

	acs, err := cfg.New()
	if err == nil {
		err = acs.Start(ctx)
	}

	if err != nil {
		_ = acs.Close()
		return nil, err
	}

	return acs, nil
}
