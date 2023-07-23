// Package darvaza-fastcgi implements a FastCGI sidecar
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"darvaza.org/sidecar/pkg/logger/zerolog"
	"darvaza.org/slog"
)

const (
	// CmdName is the name of the executable
	CmdName = "darvaza-fastcgi"

	// DefaultConfigFile is the default name for the ConfigFile
	DefaultConfigFile = CmdName + ".toml"
)

var (
	// cfg is the global Config of this tool
	cfg Config
	// cfgFile stores the config-file from command line
	cfgFile string
	// log is the global logger of this tool
	log slog.Logger
)

var rootCmd = &cobra.Command{
	Use:   CmdName,
	Short: "Darvaza FastCGI sidecar",
}

// fatal is a convenience wrapper for slog.Logger.Fatal().Print()
func fatal(err error, msg string, args ...any) {
	l := log.Fatal()
	if err != nil {
		l = l.WithField(slog.ErrorFieldName, err)
	}
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	l.Print(msg)

	panic("unreachable")
}

// main invokes cobra
func main() {
	if err := rootCmd.Execute(); err != nil {
		fatal(err, "")
	}
}

// cobraInit loads the config-file before the
// commands process their flags and arguments
func cobraInit() {
	if cfgFile != "" {
		err := cfg.ReadInFile(cfgFile)
		if err == nil {
			return
		}

		if !os.IsNotExist(err) || cfgFile != DefaultConfigFile {
			fatal(err, "failed processing %q", cfgFile)
		}

		// missing DefaultConfig, ignore
	}

	// didn't load, apply defaults
	if err := cfg.Prepare(); err != nil {
		fatal(err, "failed to set config defaults")
	}
}

// init initialises the global logger at Info level, and config-file loading
func init() {
	log = zerolog.New(nil, slog.Debug)

	// root level flags
	pflags := rootCmd.PersistentFlags()
	pflags.StringVarP(&cfgFile, "config-file", "f", DefaultConfigFile, "config file (TOML format)")

	// load config-file before the rest of the cobra commands
	cobra.OnInitialize(cobraInit)
}
