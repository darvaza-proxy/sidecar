package service

//go:generate ./args.sh

import "github.com/spf13/cobra"

// NoArgs returns an error if any args are included
func NoArgs(cmd *cobra.Command, args []string) error {
	return cobra.NoArgs(cmd, args)
}

// OnlyValidArgs returns an error if there are any positional argument
// that are not in [Config.ValidArgs]
func OnlyValidArgs(cmd *cobra.Command, args []string) error {
	return cobra.OnlyValidArgs(cmd, args)
}

// ArbitraryArgs never returns an error
func ArbitraryArgs(cmd *cobra.Command, args []string) error {
	return cobra.ArbitraryArgs(cmd, args)
}

// MaximumNArgs returns an error if there are more than N args
func MaximumNArgs(n int) cobra.PositionalArgs {
	return cobra.MaximumNArgs(n)
}

// MinimumNArgs returns an error if there are fewer than N args
func MinimumNArgs(n int) cobra.PositionalArgs {
	return cobra.MinimumNArgs(n)
}

// ExactArgs returns an error if there are not N args
func ExactArgs(n int) cobra.PositionalArgs {
	return cobra.ExactArgs(n)
}
