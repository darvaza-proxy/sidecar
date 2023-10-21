package config

import (
	"errors"
	"io"
	"path/filepath"
	"strings"

	"darvaza.org/core"
)

var (
	// ErrUnknownFormat indicates we failed to identify the format of the file
	ErrUnknownFormat = errors.New("file format not identified")
)

// A Decoder uses text to populate a data structure.
type Decoder interface {
	Decode(data string, v any) error
}

// Encoder serializes a data structure
type Encoder interface {
	SetIndent(indent string) bool
	WriteTo(v any, w io.Writer) (int64, error)
}

// NewDecoderByFilename uses the file extension
// to determine the decoder.
func NewDecoderByFilename(filename string) (Decoder, error) {
	ext := filepath.Ext(filename)
	if ext != "" {
		ext = ext[1:]
	}

	return NewDecoder(strings.ToLower(ext))
}

// NewDecoder returns the decoder associated with
// a name or extension.
func NewDecoder(name string) (Decoder, error) {
	var dec Decoder
	var key string

	if alias, ok := registryAlias[name]; ok {
		key = alias
	} else {
		key = name
	}

	if r, ok := registry[key]; ok {
		if f := r.NewDecoder; f != nil {
			dec = f()
		}
	}

	if dec == nil {
		err := core.Wrap(ErrUnknownFormat, name)
		return nil, err
	}

	return dec, nil
}

// NewEncoder returns a encoder for the specified format
func NewEncoder(name string) (Encoder, error) {
	var enc Encoder
	var key string

	if alias, ok := registryAlias[name]; ok {
		key = alias
	} else {
		key = name
	}

	if r, ok := registry[key]; ok {
		if f := r.NewEncoder; f != nil {
			enc = f()
		}
	}

	if enc == nil {
		err := core.Wrap(ErrUnknownFormat, name)
		return nil, err
	}

	return enc, nil
}

// Encoders returns all the formats we know to encode
func Encoders() []string {
	var out []string
	for name, r := range registry {
		if r.NewEncoder != nil {
			out = append(out, name)
		}
	}

	return out
}

type registryEntry struct {
	NewDecoder func() Decoder
	NewEncoder func() Encoder
}

var registry = make(map[string]*registryEntry)
var registryAlias = make(map[string]string)

func register(name string, dec func() Decoder, enc func() Encoder, aliases ...string) {
	if dec == nil || name == "" {
		panic("invalid registration")
	}

	registry[name] = &registryEntry{
		NewDecoder: dec,
		NewEncoder: enc,
	}

	for _, a := range aliases {
		registryAlias[a] = name
	}
}
