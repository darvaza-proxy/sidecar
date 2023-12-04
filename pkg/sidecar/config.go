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
	Logger  slog.Logger     `json:"-" yaml:"-" toml:"-"`
	Context context.Context `json:"-" yaml:"-" toml:"-"`
	Store   storage.Store   `json:"-" yaml:"-" toml:"-"`

	Name string `toml:"name" valid:"host,require"`

	Supervision SupervisionConfig
	Addresses   BindConfig `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	HTTP        HTTPConfig
}

// SupervisionConfig represents how graceful upgrades will operate
type SupervisionConfig struct {
	PIDFile         string        `yaml:"pid_file"         default:"/tmp/tableflip.pid"`
	GracefulTimeout time.Duration `yaml:"graceful_timeout" default:"5s"`
	HealthWait      time.Duration `yaml:"health_wait"      default:"1s"`
}

// BindConfig refers to the IP addresses used by a GoShop Server
type BindConfig struct {
	Interfaces []string `json:",omitempty"`
	Addresses  []string `json:",omitempty" valid:"ip"`
}

// HTTPConfig contains information for setting up the HTTP server
type HTTPConfig struct {
	Port              uint16        `yaml:"port"                default:"8443" valid:"port"`
	PortInsecure      uint16        `yaml:"insecure_port"       default:"8080" valid:"port"`
	EnableInsecure    bool          `yaml:"enable_insecure"`
	MutualTLSOnly     bool          `yaml:"mtls_only"`
	ReadTimeout       time.Duration `yaml:"read_timeout"        default:"1s"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" default:"2s"`
	WriteTimeout      time.Duration `yaml:"write_timeout"       default:"1s"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"        default:"30s"`
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
