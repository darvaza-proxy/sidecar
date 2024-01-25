package service

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	"darvaza.org/core"
	"darvaza.org/slog"
	"darvaza.org/x/config"
)

type cobraCmdE func(*cobra.Command, []string) error

// Service is an application that runs supervised by the OS
type Service struct {
	wg        sync.WaitGroup
	cancelled atomic.Value

	cancel context.CancelFunc
	ctx    context.Context

	log slog.Logger
	sys service.System
	ss  service.Service
	p   program

	args  []string
	run   cobraCmdE
	root  *cobra.Command
	serve *cobra.Command

	Config Config
}

// MustBuild creates a new service from a given root and serve
// [cobra.Command], and panics if there is a problem.
func MustBuild(rootCmd, serveCmd *cobra.Command) *Service {
	s, err := Build(rootCmd, serveCmd)
	if err != nil {
		panic(err)
	}
	return s
}

// Build creates a new service from a given root and serve
// [cobra.Command]. [service.Config] parameters can be modified
// until Execute is called or via PersistentPreRunE
// on the root command.
func Build(rootCmd, serveCmd *cobra.Command) (*Service, error) {
	s := newService()
	err := s.init(cmdRoot(rootCmd), cmdServe(serveCmd))
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) init(root, serve *cobra.Command) error {
	// adapt root actions
	if err := s.prepareRootCmd(root); err != nil {
		return err
	}

	// adapt serve actions
	run := s.prepareServeCmd(serve)

	// store
	s.run = run
	s.root = root
	s.serve = serve
	s.root.AddCommand(serve)

	// set initial service.Config values
	if err := s.prepareConfig(); err != nil {
		return err
	}

	// add extra commands
	switch {
	case s.sys == nil:
		// no service system detected
		s.initSolo()
	case s.sys.Interactive():
		// prepare interactive use integrated with
		// the service system.
		s.initInteractive()
	}

	return nil
}

func (s *Service) initInteractive() {
	// add service commands
	commands := []*cobra.Command{
		s.newStartCommand(),
		s.newStopCommand(),
		s.newRestartCommand(),
		s.newInstallCommand(),
		s.newUninstallCommand(),
		s.newStatusCommand(),
	}

	s.root.AddCommand(commands...)
}

func (s *Service) initSolo() {
	// serve and status commands
	commands := []*cobra.Command{
		s.newStatusCommand(),
	}

	s.root.AddCommand(commands...)
}

func (s *Service) prepareServeCmd(serve *cobra.Command) cobraCmdE {
	var appRun cobraCmdE

	switch {
	case serve.RunE != nil:
		appRun = serve.RunE
	case serve.Run != nil:
		run := serve.Run
		appRun = func(cmd *cobra.Command, arg []string) error {
			run(cmd, arg)
			return nil
		}
	}

	serve.Run = nil
	serve.RunE = s.runServe
	return appRun
}

func (s *Service) runServe(cmd *cobra.Command, args []string) error {
	if s.run == nil {
		return core.ErrNotImplemented
	}

	if s.Interactive() {
		var c core.Catcher
		return c.Do(func() error { return s.run(cmd, args) })
	}

	s.args = args
	return s.ss.Run()
}

func (s *Service) prepareRootCmd(root *cobra.Command) error {
	// hook setup
	appSetup := root.PersistentPreRunE
	setup := func(cmd *cobra.Command, args []string) error {
		return s.setup(cmd, args, appSetup)
	}
	root.PersistentPreRunE = setup
	return nil
}

func (s *Service) setup(cmd *cobra.Command, args []string, appSetup cobraCmdE) error {
	ctx := cmd.Context()
	if err := s.initContext(ctx); err != nil {
		return err
	}

	if appSetup != nil {
		// run app's own setup
		if err := appSetup(cmd, args); err != nil {
			return err
		}
	}

	return s.prepareService()
}

func (s *Service) prepareService() error {
	if err := s.prepareConfig(); err != nil {
		return core.Wrap(err, "Service.Config")
	}

	sc := &s.Config.Config
	ss, err := service.New(&s.p, sc)
	if err != nil {
		return err
	}

	s.ss = ss
	return nil
}

func (s *Service) prepareConfig() error {
	// Name
	name := cmdUseName(s.root, CmdName())
	if s.Config.Name == "" {
		s.Config.Name = name
	}
	if s.Config.DisplayName == "" {
		s.Config.DisplayName = name
	}
	// Description
	if s.Config.Description == "" {
		s.Config.Description = cmdDescription(s.root, "")
	}

	// TODO: validate

	return config.SetDefaults(&s.Config)
}

func (s *Service) newStartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Starts the service",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return s.ss.Start()
		},
		SilenceUsage: true,
	}
}

func (s *Service) newStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stops the service",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return s.ss.Stop()
		},
		SilenceUsage: true,
	}
}

func (s *Service) newRestartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restarts the service",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return s.ss.Restart()
		},
		SilenceUsage: true,
	}
}

func (s *Service) newInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Installs the service on the system",
		RunE: func(_ *cobra.Command, _ []string) error {
			return s.ss.Install()
		},
		SilenceUsage: true,
	}
}

func (s *Service) newUninstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstalls the service from the system",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return s.ss.Uninstall()
		},
		SilenceUsage: true,
	}
}

func (s *Service) newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Shows the current service status",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			var res string

			st, err := s.ss.Status()
			if err != nil {
				return err
			}

			labels := []string{"unknown", "running", "stopped"}
			if int(st) < len(labels) {
				res = labels[st]
			} else {
				res = fmt.Sprintf("%s(%v)", labels[0], int(st))
			}

			_, err = fmt.Printf("%s: %s\n", "Status", res)
			return err
		},
		SilenceUsage: true,
	}
}

// Interactive returns false if running under the OS service
// manager and true otherwise.
func (s *Service) Interactive() bool {
	if s.sys != nil {
		return s.sys.Interactive()
	}
	return true
}

// Platform returns a description of the system service.
func (s *Service) Platform() string {
	if s.sys != nil {
		return s.sys.String()
	}
	return ""
}

func newService() *Service {
	s := new(Service)
	s.p.s = s
	s.sys = service.ChosenSystem()
	return s
}

type program struct {
	s *Service
}

func (p *program) Start(_ service.Service) error {
	return p.s.doStart()
}

func (p *program) Stop(_ service.Service) error {
	return p.s.doStop()
}
