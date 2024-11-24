package main

import (
	"darvaza.org/sidecar/pkg/sidecar"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run Edgy instance",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		srv, err := sidecar.New(&srvConf.Config)
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
