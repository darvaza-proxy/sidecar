package config

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	"darvaza.org/x/config"
	"darvaza.org/x/config/appdir"
	"darvaza.org/x/config/expand"
)

var (
	// CmdName is the base name of default config files if none is
	// specified on the [Loader]
	CmdName = filepath.Base(os.Args[0])
	// DefaultExtensions is the list of extensions to test when trying
	// to find a config file.
	DefaultExtensions = []string{"conf", "json", "toml", "yaml", "yml"}
)

// Loader helps finding and location configuration files
type Loader[T any] struct {
	l config.Loader[T]

	Base        string
	Directories []string
	Extensions  []string
	Others      []string
}

// Last returns the filename of the last config file to be used
func (l *Loader[T]) Last() (fs.FS, string) {
	return l.l.Last()
}

// SetDefaults fills the gaps in the [Loader] config. This is not required
// as LoadFileFlag() and LoadKnownLocations() will call it
// automatically before trying to find
func (l *Loader[T]) SetDefaults() {
	if l.Base == "" {
		l.Base = CmdName
	}

	if len(l.Directories) == 0 {
		l.Directories = appdir.AllConfigDir(CmdName)
	}

	if len(l.Extensions) == 0 {
		l.Extensions = DefaultExtensions
	}

	l.l.NewDecoder = NewDecoderFactory[T](nil)
}

// NewFromFile loads a config file by name, followed by initialization, filling gaps,
// and validation.
func (l *Loader[T]) NewFromFile(configFile string, options ...config.Option[T]) (*T, error) {
	l.syncOptions(options)
	return l.l.NewFromFileOS(configFile)
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
	l.SetDefaults()

	files, err := config.Join(l.Directories, l.Base, l.Extensions)
	if err != nil {
		return nil, err
	}
	files = append(files, l.Others...)

	l.syncOptions(options)
	v, err := l.l.NewFromFileOS(files...)
	switch {
	case err == nil:
		return v, nil
	case os.IsNotExist(err):
		return l.l.New()
	default:
		return nil, err
	}
}

func (l *Loader[T]) syncOptions(options []config.Option[T]) {
	l.l.Options = options
}

// NewDecoderFactory returns a [config.Decoder] factory to use with [config.Loader],
// using our registered decoders.
func NewDecoderFactory[T any](getenv func(string) string) func(string) (config.Decoder[T], error) {
	return func(filename string) (config.Decoder[T], error) {
		dec, _ := NewDecoderByFilename(filename)
		if dec == nil {
			// fallback to autodetect
			dec, _ = NewDecoder("auto")
		}

		if dec == nil {
			return nil, config.NewPathError(filename, "decode", ErrUnknownFormat)
		}
		fn := newLoaderDecoder[T](dec, getenv)
		return fn, nil
	}
}

func newLoaderDecoder[T any](dec Decoder, getenv func(string) string) config.DecoderFunc[T] {
	return func(filename string, data []byte) (*T, error) {
		text, err := expand.FromBytes(data, getenv)
		if err != nil {
			return nil, config.NewPathError(filename, "decode", err)
		}

		v := new(T)
		err = dec.Decode(text, v)
		if err != nil {
			return nil, config.NewPathError(filename, "decode", err)
		}

		return v, nil
	}
}
