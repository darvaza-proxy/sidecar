// Package store provides a TLS certificate manager for sidecars
package store

import (
	"context"
	"os"
	"time"

	"darvaza.org/core"
	"darvaza.org/slog"
	"darvaza.org/x/tls"
	"darvaza.org/x/tls/store"
	"darvaza.org/x/tls/store/basic"
)

const (
	// DefaultRootsFile ...
	DefaultRootsFile = "caroot.pem"
	// DefaultPrivateKeyFile ...
	DefaultPrivateKeyFile = "key.pem"
	// DefaultCertFile ...
	DefaultCertFile = "cert.pem"

	// ConfigNewTimeout indicates how low we allow for certificates and keys to be
	// loaded.
	ConfigNewTimeout = 1 * time.Second
)

// Config contains information for setting up TLS clients and server
type Config struct {
	Key   string `default:"key.pem"`
	Cert  string `default:"cert.pem"`
	Roots string `default:"caroot.pem"`
}

// New creates a simple [storage.Store] from the config
func (cfg *Config) New(logger slog.Logger) (tls.Store, error) {
	tio := time.Now().Add(ConfigNewTimeout)
	ctx, cancel := context.WithDeadline(context.Background(), tio)
	defer cancel()

	s := basic.New()
	if err := cfg.apply(ctx, logger, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (cfg *Config) apply(ctx context.Context, logger slog.Logger, out tls.StoreX509Writer) error {
	var errs core.CompoundError

	if err := cfg.applyCACerts(ctx, out, logger, DefaultRootsFile); err != nil {
		errs.AppendError(err)
	}

	if err := cfg.applyPrivateKey(ctx, out, logger, DefaultPrivateKeyFile); err != nil {
		errs.AppendError(err)
	}

	if err := cfg.applyCerts(ctx, out, logger, DefaultCertFile); err != nil {
		errs.AppendError(err)
	}

	return errs.AsError()
}

func (cfg *Config) applyCACerts(ctx context.Context, out tls.StoreX509Writer,
	logger slog.Logger, defaultFile string) error {
	//
	value := cfg.Roots
	if value == "" {
		return nil
	}

	sc := &store.Config{
		Logger: logger,
	}

	_, err := sc.AddCACerts(ctx, out, value)
	if os.IsNotExist(err) && value == defaultFile {
		err = nil
	}

	return err
}

func (cfg *Config) applyPrivateKey(ctx context.Context, out tls.StoreX509Writer,
	logger slog.Logger, defaultFile string) error {
	//
	value := cfg.Key
	if value == "" {
		return nil
	}

	sc := &store.Config{
		Logger: logger,
	}

	err := sc.AddPrivateKey(ctx, out, value)
	if os.IsNotExist(err) && value == defaultFile {
		err = nil
	}

	return err
}

func (cfg *Config) applyCerts(ctx context.Context, out tls.StoreX509Writer,
	logger slog.Logger, defaultFile string) error {
	//
	value := cfg.Cert
	if value == "" {
		return nil
	}

	sc := &store.Config{
		Logger: logger,
	}

	err := sc.AddCert(ctx, out, value)
	if os.IsNotExist(err) && value == defaultFile {
		err = nil
	}
	return err
}
