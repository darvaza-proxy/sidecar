package main

import (
	"context"

	"github.com/spf13/pflag"

	"darvaza.org/sidecar/pkg/config"
	"darvaza.org/sidecar/pkg/sidecar"
	"darvaza.org/slog"
)

// Config ...
type Config struct {
	sidecar.Config

	TLS ACMEConfig
}

const (
	configFileFlag      = "config"
	configFileShortFlag = "f"
	configFileDefault   = CmdName + ".{conf,json,toml,yaml}"
)

var confLoader = config.Loader[Config]{
	Base:        CmdName,
	Directories: []string{"."},
	Extensions:  []string{"conf", "json", "toml", "yaml"},
}

func getConfig(ctx context.Context, flags *pflag.FlagSet) (*Config, error) {
	log := getLogger(ctx, flags)
	init := func(cfg *Config) error {
		cfg.Context = ctx
		cfg.Logger = log
		return nil
	}

	flag := flags.Lookup(configFileFlag)
	cfg, err := confLoader.NewFromFlag(flag, init)
	if err != nil {
		log.Error().WithField(slog.ErrorFieldName, err).Print("LoadConfigFile")
		return nil, err
	}

	if _, s := confLoader.Last(); s != "" {
		log.Info().WithField("filename", s).Print("config loaded")
	}

	return cfg, nil
}

func init() {
	pFlags := rootCmd.PersistentFlags()
	pFlags.StringP(configFileFlag, configFileShortFlag, configFileDefault, "config file to use")
}
