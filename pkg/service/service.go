package service

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

// Service is an application that runs supervised by the OS
type Service struct {
	wg         sync.WaitGroup
	cancelOnce sync.Once

	cancel context.CancelFunc
	ctx    context.Context

	d   time.Duration
	err atomic.Value
	ss  service.Service
	p   program

	cmd     cobra.Command
	prepare func(context.Context, *cobra.Command, []string) error
	run     func(context.Context, *cobra.Command, []string) error

	initializers []func()
	finalizers   []func()
}

// Must creates a new service and panics if there is a problem
func Must(cfg *Config) *Service {
	s, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return s
}

// New creates a new service
func New(cfg *Config) (*Service, error) {
	if cfg == nil {
		cfg = new(Config)
	}

	if err := cfg.SetDefaults(); err != nil {
		return nil, err
	}

	s := new(Service)

	sc := &service.Config{
		Name:        cfg.Name,
		DisplayName: cfg.DisplayName,
		Description: cfg.Description,
	}

	ss, err := service.New(&s.p, sc)
	if err != nil {
		return nil, err
	}

	if err := s.init(ss, cfg); err != nil {
		return nil, err
	}

	// add service commands
	commands := []*cobra.Command{
		s.newStartCommand(),
		s.newStopCommand(),
		s.newRestartCommand(),
		s.newInstallCommand(),
		s.newUninstallCommand(),
		s.newStatusCommand(),
	}

	s.cmd.AddCommand(commands...)
	return s, nil
}

func (s *Service) init(ss service.Service, cfg *Config) error {
	ctx, cancel := context.WithCancel(cfg.Context)

	// populate
	*s = Service{
		cancel: cancel,
		ctx:    ctx,

		d:  cfg.SanityDelay,
		ss: ss,
		p: program{
			s: s,
		},

		prepare: cfg.Prepare,
		run:     cfg.Run,

		cmd: cobra.Command{
			Use:     cfg.Name,
			Short:   cfg.Short,
			Version: cfg.Version,

			Args:      cfg.ValidateArgs,
			ValidArgs: cfg.ValidArgs,

			RunE: func(cmd *cobra.Command, args []string) error {
				if err := s.prepare(ctx, cmd, args); err != nil {
					return err
				}

				return s.run(ctx, cmd, args)
			},
		},
	}

	return nil
}

func (s *Service) newStartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "starts the service",
		RunE: func(_ *cobra.Command, _ []string) error {
			return s.ss.Start()
		},
	}
}

func (s *Service) newStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "stops the service",
		RunE: func(_ *cobra.Command, _ []string) error {
			return s.ss.Stop()
		},
	}
}

func (s *Service) newRestartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "restarts the service",
		RunE: func(_ *cobra.Command, _ []string) error {
			return s.ss.Restart()
		},
	}
}

func (s *Service) newInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "installs the service on the system",
		RunE: func(_ *cobra.Command, _ []string) error {
			return s.ss.Install()
		},
	}
}

func (s *Service) newUninstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "uninstalls the service from the system",
		RunE: func(_ *cobra.Command, _ []string) error {
			return s.ss.Uninstall()
		},
	}
}

func (s *Service) newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "shows the current service status",
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
	}
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
