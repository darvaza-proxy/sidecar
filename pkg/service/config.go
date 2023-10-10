package service

import (
	"context"
	"time"

	"darvaza.org/core"
	"github.com/spf13/cobra"
)

const (
	// DefaultSanityDelay indicates how long we wait for
	// the run command to fail before claiming a
	// successful start
	DefaultSanityDelay = 2 * time.Second
)

// Config describes the Service we are assembling
type Config struct {
	Name        string
	DisplayName string
	Description string
	Short       string
	Version     string

	Context     context.Context
	SanityDelay time.Duration

	Prepare func(context.Context, *cobra.Command, []string) error
	Run     func(context.Context, *cobra.Command, []string) error

	ValidArgs    []string
	ValidateArgs cobra.PositionalArgs
}

// SetDefaults fills gaps in the config
func (cfg *Config) SetDefaults() error {
	if cfg.Name == "" {
		cfg.Name = CmdName()
	}

	if cfg.DisplayName == "" {
		cfg.DisplayName = cfg.Name
	}

	if cfg.Short == "" {
		cfg.Short = cfg.Name + "runs the command as a service."
	}

	if cfg.Context == nil {
		cfg.Context = context.Background()
	}

	if cfg.SanityDelay == 0 {
		cfg.SanityDelay = DefaultSanityDelay
	}

	cfg.setDefaultEntrypoints()
	return nil
}

func (cfg *Config) setDefaultEntrypoints() {
	switch {
	case cfg.Run == nil:
		// no command
		cfg.Run = notImplementedCommand

		if cfg.Prepare == nil {
			// abort early
			cfg.Prepare = notImplementedCommand
		}
	case cfg.Prepare == nil:
		// no preparation needed
		cfg.Prepare = noOperationCommand
	}

	switch {
	case cfg.ValidateArgs != nil:
		// custom
	case len(cfg.ValidArgs) == 0:
		// no args
		cfg.ValidateArgs = cobra.NoArgs
	default:
		// only specified args
		cfg.ValidateArgs = cobra.OnlyValidArgs
	}
}

func notImplementedCommand(_ context.Context, _ *cobra.Command, _ []string) error {
	return core.ErrNotImplemented
}

func noOperationCommand(_ context.Context, _ *cobra.Command, _ []string) error {
	return nil
}
