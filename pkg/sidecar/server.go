// Package sidecar implements the common server side of all Darvaza sidecars
package sidecar

import (
	"context"
	"sync/atomic"

	"darvaza.org/core"
)

// Server is the HTTP Server of the sidecar
type Server struct {
	cfg       Config
	ctx       context.Context
	cancel    context.CancelFunc
	cancelled atomic.Bool
	err       atomic.Value
	wg        core.WaitGroup
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

	// and continue
	return cfg.newServer()
}

// New creates a new Server from the Config
func (cfg *Config) New() (*Server, error) {
	return New(cfg)
}

func (cfg *Config) newServer() (*Server, error) {
	ctx, cancel := context.WithCancel(cfg.Context)

	srv := &Server{
		cfg:    *cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	srv.wg.OnError(srv.onWorkerError)

	return srv, nil
}
