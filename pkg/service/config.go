package service

import (
	"time"

	"github.com/kardianos/service"
)

// Config describes the Service we are assembling
type Config struct {
	service.Config

	// SanityDelay indicates how long we wait for
	// the run command to fail before claiming a
	// successful start
	SanityDelay time.Duration `default:"2s"`

	// GracefulTimeout indicates how long we wait for
	// the shutdown process to take before forcing an exit.
	GracefulTimeout time.Duration `default:"10s"`
}

// SetDefaults fills gaps in the Config
func (cfg *Config) SetDefaults() {
	if cfg.Name == "" {
		cfg.Name = CmdName()
	}

	cfg.setOSDefaults()
}
