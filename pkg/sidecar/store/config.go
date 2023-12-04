// Package store provides a TLS certificate manager for sidecars
package store

import (
	"darvaza.org/darvaza/shared/storage"
	"darvaza.org/darvaza/shared/storage/simple"
	"darvaza.org/slog"
)

// Config contains information for setting up TLS clients and server
type Config struct {
	Key   string `default:"key.pem"`
	Cert  string `default:"cert.pem"`
	Roots string `default:"caroot.pem"`
}

// New creates a simple [storage.Store] from the config
func (cfg *Config) New(logger slog.Logger) (storage.Store, error) {
	sc := &simple.Config{
		Logger: logger,
	}

	return sc.New(cfg.Key, cfg.Cert, cfg.Roots)
}
