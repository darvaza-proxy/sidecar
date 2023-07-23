package main

import (
	"net/http"

	"darvaza.org/sidecar/pkg/config/flags"
	"darvaza.org/sidecar/pkg/config/flags/cobra"
)

// Command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "runs the sidecar",
	PreRun: func(cmd *cobra.Command, _ []string) {
		flags.GetMapper(cmd.Flags()).Parse()
	},
	RunE: func(_ *cobra.Command, _ []string) error {
		var r http.Handler

		// Logger
		cfg.Server.Logger = log

		// Prepare server
		srv, err := cfg.Server.New()
		if err != nil {
			return err
		}

		return srv.ListenAndServe(r)
	},
}

func init() {
	// Flags
	cobra.NewMapper(serveCmd.Flags()).
		VarP(&cfg.Server.HTTP.Port, "port", 'p', "HTTPS Port").
		Var(&cfg.Server.Supervision.PIDFile, "pid-file", "Path to PID file").
		VarP(&cfg.Server.Supervision.GracefulTimeout, "graceful-timeout", 't',
			"Maximum time to wait for in-flight requests")

	rootCmd.AddCommand(serveCmd)
}
