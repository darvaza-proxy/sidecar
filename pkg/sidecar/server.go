// Package sidecar implements the common engine of all Darvaza sidecars
package sidecar

import (
	"context"
	"sync/atomic"

	"darvaza.org/core"
	"darvaza.org/darvaza/agent/httpserver"
	"darvaza.org/darvaza/shared/storage"
)

// Server is the HTTP Server of the sidecar
type Server struct {
	cfg       Config
	ctx       context.Context
	cancel    context.CancelFunc
	cancelled atomic.Bool
	err       atomic.Value
	wg        core.WaitGroup

	tls storage.Store
	hs  *httpserver.Server
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
		s, err = newTLSStore(cfg)
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
	ctx, cancel := context.WithCancel(cfg.Context)

	srv := &Server{
		cfg:    *cfg,
		ctx:    ctx,
		cancel: cancel,
		tls:    s,
	}

	srv.wg.OnError(srv.onWorkerError)

	hsc := srv.newHTTPServerConfig()
	hs, err := hsc.New()
	if err != nil {
		return nil, err
	}
	srv.hs = hs

	return srv, nil
}
