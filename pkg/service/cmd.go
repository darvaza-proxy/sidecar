package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"darvaza.org/core"
)

// CmdName returns the arg[0] of this execution
func CmdName() string {
	return filepath.Base(os.Args[0])
}

// cmdUseName attempts to extract a Service.Name
// from a cobra.Command
func cmdUseName(cmd *cobra.Command, defaultValue string) string {
	if cmd != nil {
		name := cmd.Use
		idx := strings.IndexFunc(name, unicode.IsSpace)
		if idx >= 0 {
			name = name[:idx]
		}
		if name != "" {
			return name
		}
	}

	return defaultValue
}

// cmdDescription attempts to extract a Service.Description
// from a cobra.Command.
func cmdDescription(cmd *cobra.Command, defaultValue string) string {
	if cmd != nil {
		return core.Coalesce(cmd.Long, cmd.Short, defaultValue)
	}

	return defaultValue
}

// cmdRoot fills the gaps in the root command.
func cmdRoot(cmd *cobra.Command) *cobra.Command {
	cmdName := CmdName()

	cmd = core.Coalesce(cmd, &cobra.Command{
		SilenceErrors: true,
		SilenceUsage:  true,
	})

	cmd.Use = core.Coalesce(cmd.Use, cmdName)

	if cmd.Short == "" {
		// TODO: first letter to uppercase
		name := cmdUseName(cmd, cmdName)
		cmd.Short = name + " service."
	}

	return cmd
}

// cmdServe fills the gaps in the serve command.
func cmdServe(cmd *cobra.Command) *cobra.Command {
	cmd = core.Coalesce(cmd, &cobra.Command{
		SilenceErrors: true,
		SilenceUsage:  true,
	})

	cmd.Use = core.Coalesce(cmd.Use, "serve")
	cmd.Short = core.Coalesce(cmd.Short, "Runs the service.")

	if cmd.Run == nil && cmd.RunE == nil {
		cmd.Args = core.Coalesce(cmd.Args, cobra.NoArgs)

		cmd.RunE = func(_ *cobra.Command, _ []string) error {
			return core.ErrNotImplemented
		}
	}

	return cmd
}

// Execute either runs via cobra or the service invoked process
func (s *Service) Execute() error {
	return s.ExecuteContext(context.Background())
}

// ExecuteContext either runs via cobra or the service invoked process
func (s *Service) ExecuteContext(ctx context.Context) error {
	ctx, err := s.NewContext(ctx)
	if err != nil {
		return err
	}

	return s.root.ExecuteContext(ctx)
}

// Command returns the root [cobra.Command]
func (s *Service) Command() *cobra.Command {
	return s.root
}

// AddCommand adds an extra interactive command to the tool
func (s *Service) AddCommand(cmd *cobra.Command) {
	s.root.AddCommand(cmd)
}

// PersistentFlags returns the flags also passed to sub-commands
func (s *Service) PersistentFlags() *pflag.FlagSet {
	return s.root.PersistentFlags()
}
