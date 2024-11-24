package main

import (
	"darvaza.org/core"
	"darvaza.org/sidecar/pkg/sidecar"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run Edgy instance",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg := srvConf

		// TLS Store
		acme, err := cfg.TLS.New(cfg.Logger)
		if err != nil {
			return core.Wrap(err, "failed to create autocert store")
		}

		if err := acme.Start(cfg.Context); err != nil {
			_ = acme.Close()
			return core.Wrap(err, "failed to start autocert worker")
		}
		defer closeAll(acme)

		cfg.Config.Store = acme
		srv, err := sidecar.New(&cfg.Config)
		if err != nil {
			return err
		}

		return srv.ListenAndServe(nil)
	},
}

// WantsSyslog tells if the `--syslog` flag was passed
// to use the system logger in interactive mode.
func WantsSyslog(flags *pflag.FlagSet) bool {
	v, _ := flags.GetBool(syslogFlag)
	return v
}

const syslogFlag = "syslog"

func init() {
	flags := serveCmd.Flags()
	flags.Bool(syslogFlag, false, "use syslog when running manually")
}
