package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	"darvaza.org/x/config"
)

var (
	// CmdName is the base name of default config files if none is
	// specified on the [Loader]
	CmdName = filepath.Base(os.Args[0])
	// DefaultDirectories is the list of directories to test when trying
	// to find a config file.
	DefaultDirectories = []string{".", "/etc", "/etc/" + CmdName}
	// DefaultExtensions is the list of extensions to test when trying
	// to find a config file.
	DefaultExtensions = []string{"conf", "json", "toml", "yaml", "yml"}
)

// Loader helps finding and location configuration files
type Loader[T any] struct {
	last string

	Base        string
	Directories []string
	Extensions  []string
	Others      []string
}

// Last returns the filename of the last config file to be used
func (l *Loader[T]) Last() string {
	return l.last
}

// SetDefaults fills the gaps in the [Loader] config. This is not required
// as LoadFileFlag() and LoadKnownLocations() will call it
// automatically before trying to find
func (l *Loader[T]) SetDefaults() {
	if l.Base == "" {
		l.Base = CmdName
	}

	if len(l.Directories) == 0 {
		l.Directories = DefaultDirectories
	}

	if len(l.Extensions) == 0 {
		l.Extensions = DefaultExtensions
	}
}

// NewFromFile loads a config file by name, followed by initialization, filling gaps,
// and validation.
func (l *Loader[T]) NewFromFile(configFile string, options ...config.Option[T]) (*T, error) {
	l.last = configFile

	cfg := new(T)
	err := LoadFile(configFile, cfg, options...)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// NewFromFlag uses a cobra Flag for the config file name. if not specifies
// the known locations will be tried, and if all fails, it will create a
// default one.
func (l *Loader[T]) NewFromFlag(flag *pflag.Flag, options ...config.Option[T]) (*T, error) {
	if flag.Changed {
		// given
		configFile := flag.Value.String()
		return l.NewFromFile(configFile, options...)
	}

	return l.NewFromKnownLocations(options...)
}

// NewFromKnownLocations scans locations specified in the [Loader] for a config file,
// if not possible it will create a default one.
func (l *Loader[T]) NewFromKnownLocations(options ...config.Option[T]) (*T, error) {
	var files []string

	l.SetDefaults()

	files = append(files, l.Others...)
	for _, dir := range l.Directories {
		for _, ext := range l.Extensions {
			fn := fmt.Sprintf("%s/%s.%s", dir, l.Base, ext)
			files = append(files, fn)
		}
	}

	for _, fn := range files {
		cfg, err := l.NewFromFile(fn, options...)
		switch {
		case cfg != nil:
			// good file found
			return cfg, nil
		case err != nil && !os.IsNotExist(err):
			// bad file found
			return nil, err
		}
	}

	// make one fresh
	l.last = ""
	return New(options...)
}
