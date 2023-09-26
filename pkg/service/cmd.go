package service

import (
	"os"
	"path/filepath"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CmdName returns the arg[0] of this execution
func CmdName() string {
	return filepath.Base(os.Args[0])
}

// Execute either runs via cobra or the service invoked process
func (s *Service) Execute() error {
	if service.Interactive() {
		// manual
		return s.cmd.Execute()
	}

	// service
	if err := s.cmd.ParseFlags(os.Args[1:]); err != nil {
		return err
	}

	return s.ss.Run()
}

func (s *Service) preRun() {
	for _, x := range s.initializers {
		x()
	}
}

func (s *Service) postRun() {
	for _, x := range s.finalizers {
		x()
	}
}

// Command returns the root [cobra.Command]
func (s *Service) Command() *cobra.Command {
	return &s.cmd
}

// AddCommand adds an extra interactive command to the tool
func (s *Service) AddCommand(cmd *cobra.Command) {
	s.cmd.AddCommand(cmd)
}

// PersistentFlags returns the flags also passed to sub-commands
func (s *Service) PersistentFlags() *pflag.FlagSet {
	return s.cmd.PersistentFlags()
}

// LocalFlags returns the flags only used for the root command
func (s *Service) LocalFlags() *pflag.FlagSet {
	return s.cmd.LocalFlags()
}

// OnInitialize sets functions to be executed right after
// parsing the arguments.
func (s *Service) OnInitialize(preRun ...func()) {
	cobra.OnInitialize(preRun...)
	s.initializers = append(s.initializers, preRun...)
}

// OnFinalize sets functions to be executed after running
// the commands.
func (s *Service) OnFinalize(postRun ...func()) {
	cobra.OnFinalize(postRun...)
	s.finalizers = append(s.finalizers, postRun...)
}
