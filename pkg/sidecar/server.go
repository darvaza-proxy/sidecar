// Package sidecar implements the common engine of all Darvaza sidecars
package sidecar

import (
	"darvaza.org/core"
	"darvaza.org/darvaza/shared/storage"
	"darvaza.org/sidecar/pkg/sidecar/dnsserver"
	"darvaza.org/sidecar/pkg/sidecar/httpserver"
)

// Server is the HTTP Server of the sidecar
type Server struct {
	cfg Config
	eg  core.ErrGroup

	tls storage.Store
	hs  *httpserver.Server
	ds  *dnsserver.Server
}

// New creates a new HTTP [Server] using the given [Config]
func New(cfg *Config) (*Server, error) {
	var err error

	// prepare Config
	if cfg == nil {
		cfg = &Config{}

		if err = cfg.SetDefaults(); err != nil {
			return nil, err
		}
	}

	if err = cfg.Validate(); err != nil {
		return nil, err
	}

	// TLS
	s := cfg.Store
	if s == nil {
		s, err = newFallbackTLSStore()
		if err != nil {
			return nil, err
		}
	}

	// and continue
	return cfg.newServer(s)
}

// New creates a new Server from the Config
func (cfg *Config) New() (*Server, error) {
	return New(cfg)
}

// NewWithStore creates a new server using the given config and
// a prebuilt tls Store
func (cfg *Config) NewWithStore(s storage.Store) (*Server, error) {
	cfg.Store = s
	return New(cfg)
}

func (cfg *Config) newServer(s storage.Store) (*Server, error) {
	srv := &Server{
		cfg: *cfg,
		eg: core.ErrGroup{
			Parent: cfg.Context,
		},
		tls: s,
	}

	srv.eg.OnError(srv.onWorkerError)

	if err := srv.init(); err != nil {
		return nil, err
	}

	return srv, nil
}

func (srv *Server) init() error {
	for _, fn := range []func() error{
		srv.initAddresses,
		srv.initHTTPServer,
		srv.initDNSServer,
	} {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}
