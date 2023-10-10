package sidecar

import (
	"context"
	"fmt"
	"time"

	"darvaza.org/darvaza/shared/config"
	"darvaza.org/darvaza/shared/storage"
	"darvaza.org/slog"
	"darvaza.org/slog/handlers/discard"
)

// Config represents the generic server configuration for Darvaza sidecars
type Config struct {
	Logger  slog.Logger     `toml:"-"`
	Context context.Context `toml:"-"`
	Store   storage.Store   `toml:"-"`

	Name string `toml:"name" valid:"host,require"`

	Supervision SupervisionConfig `toml:"run"`
	Addresses   BindConfig        `toml:",omitempty"`
	HTTP        HTTPConfig        `toml:"http"`
}

// SupervisionConfig represents how graceful upgrades will operate
type SupervisionConfig struct {
	PIDFile         string        `toml:"pid_file"         default:"/tmp/tableflip.pid"`
	GracefulTimeout time.Duration `toml:"graceful_timeout" default:"5s"`
	HealthWait      time.Duration `toml:"health_wait"      default:"1s"`
}

// BindConfig refers to the IP addresses used by a GoShop Server
type BindConfig struct {
	Interfaces []string `toml:"interfaces"`
	Addresses  []string `toml:"addresses" valid:"ip"`
}

// HTTPConfig contains information for setting up the HTTP server
type HTTPConfig struct {
	Port              uint16        `toml:"port"                default:"8443" valid:"port"`
	PortInsecure      uint16        `toml:"insecure_port"       default:"8080" valid:"port"`
	EnableInsecure    bool          `toml:"enable_insecure"`
	MutualTLSOnly     bool          `toml:"mtls_only"`
	ReadTimeout       time.Duration `toml:"read_timeout"        default:"1s"`
	ReadHeaderTimeout time.Duration `toml:"read_header_timeout" default:"2s"`
	WriteTimeout      time.Duration `toml:"write_timeout"       default:"1s"`
	IdleTimeout       time.Duration `toml:"idle_timeout"        default:"30s"`
}

// SetDefaults fills the gaps in the Config
func (cfg *Config) SetDefaults() error {
	if cfg.Logger == nil {
		cfg.Logger = discard.New()
	}

	if cfg.Context == nil {
		cfg.Context = context.Background()
	}

	return config.SetDefaults(cfg)
}

// Validate tells if the configuration is worth a try
func (cfg *Config) Validate() error {
	err := config.Validate(cfg)
	if err != nil {
		return err
	}

	// context.Background is *0 so valid:",required" fails
	if cfg.Context == nil {
		return fmt.Errorf("%s: %s", "Context", "can not be nil")
	}

	return nil
}
