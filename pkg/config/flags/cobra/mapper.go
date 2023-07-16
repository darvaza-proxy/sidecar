package cobra

import (
	"github.com/spf13/pflag"

	"darvaza.org/sidecar/pkg/config/flags"
)

var (
	_ flags.Mapper = (*Mapper)(nil)
)

// Mapper is a flags.Mapper for spf13's pflag
type Mapper struct {
	flags *pflag.FlagSet
}

// Parse applies flags to their mapped config fields
func (*Mapper) Parse() error {
	return nil
}

// Var maps a pflag longOption to a config field
func (m *Mapper) Var(_ any, _ string, _ string) flags.Mapper {
	return m
}

// VarP maps a pflag longOption and shortOption to a config field
func (m *Mapper) VarP(_ any, _ string, _ rune, _ string) flags.Mapper {
	return m
}

// NewMapper creates and registers a flags.Mapper for the given [*pflag.Flagset]
func NewMapper(flagSet *pflag.FlagSet) flags.Mapper {
	if flagSet != nil {
		m := &Mapper{flags: flagSet}
		flags.PutMapper(flagSet, m)
		return m
	}
	return nil
}
