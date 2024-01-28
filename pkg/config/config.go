// Package config contains helpers for loading and dumping TOML config files
package config

import (
	"darvaza.org/x/config"
	"darvaza.org/x/config/expand"
)

// LoadFile loads a config file by name, expanding environment variables,
// filling gaps and validating its content.
// LoadFile uses the file extension to determine the format.
func LoadFile[T any](filename string, v *T,
	options ...config.Option[T]) error {
	//
	data, err := expand.FromFile(filename, nil)
	if err != nil {
		// failed to read
		return err
	}

	dec, _ := NewDecoderByFilename(filename)
	if dec == nil {
		// fallback to autodetect
		dec, _ = NewDecoder("auto")
	}

	if dec == nil {
		return config.NewPathError(filename, "decode", ErrUnknownFormat)
	}

	if err := dec.Decode(data, v); err != nil {
		// failed to decode
		return config.NewPathError(filename, "decode", err)
	}

	return loadFileDecoded[T](filename, v, options)
}

func loadFileDecoded[T any](filename string, v *T, options []config.Option[T]) error {
	for _, opt := range options {
		if err := opt(v); err != nil {
			return config.NewPathError(filename, "init", err)
		}
	}

	if err := config.Prepare(v); err != nil {
		// failed to validate
		return config.NewPathError(filename, "validate", err)
	}

	// success
	return nil
}

// New creates a new config, applying the initialization functions,
// filling gaps and validating its content.
func New[T any](options ...config.Option[T]) (*T, error) {
	v := new(T)
	if err := loadFileDecoded("", v, options); err != nil {
		return nil, err
	}
	return v, nil
}

// Prepare fills any gap in the object and validates its content.
func Prepare(v any) error {
	return config.Prepare(v)
}
