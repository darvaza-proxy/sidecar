package main

import (
	"io"
	"os"

	"darvaza.org/sidecar/pkg/config"
	"darvaza.org/sidecar/pkg/config/flags/cobra"
	"darvaza.org/sidecar/pkg/sidecar"
)

// Config is the configuration structure of this sidecar
type Config struct {
	Server sidecar.Config `toml:"server"`
}

// ReadInFile loads the sidecar configuration from a TOML file
// by name, expanding environment variables, filling gaps and
// validating its content. On error the object isn't touched.
func (cfg *Config) ReadInFile(filename string) error {
	var c Config

	err := config.LoadFile(filename, &c)
	if err != nil {
		return err
	}

	*cfg = c
	return nil
}

// WriteTo writes out the Config encoded as TOML
func (cfg *Config) WriteTo(w io.Writer) (int64, error) {
	return config.WriteTo(cfg, w)
}

// Prepare fills any gap in the Config and validates its content
func (cfg *Config) Prepare() error {
	return config.Prepare(cfg)
}

// Command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump shows the loaded config",
	RunE: func(_ *cobra.Command, _ []string) error {
		if _, err := cfg.WriteTo(os.Stdout); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}
