// Package config contains helpers for loading and dumping TOML config files
package config

import (
	"os"

	"darvaza.org/darvaza/shared/config"
	"darvaza.org/darvaza/shared/config/expand"
)

// LoadFile loads a config file by name, expanding environment variables,
// filling gaps and validating its content.
// LoadFile uses the file extension to determine the format.
func LoadFile[T any](filename string, v *T,
	options ...func(*T) error) error {
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
		return ErrUnknownFormat
	}

	if err := dec.Decode(data, v); err != nil {
		// failed to decode
		err = &os.PathError{
			Path: filename,
			Op:   "decode",
			Err:  err,
		}
		return err
	}

	return loadFileDecoded[T](filename, v, options)
}

func loadFileDecoded[T any](filename string, v *T, options []func(*T) error) error {
	for _, opt := range options {
		if err := opt(v); err != nil {
			err = &os.PathError{
				Path: filename,
				Op:   "init",
				Err:  err,
			}
			return err
		}
	}

	if err := config.Prepare(v); err != nil {
		// failed to validate
		err = &os.PathError{
			Path: filename,
			Op:   "validate",
			Err:  err,
		}
		return err
	}

	// success
	return nil
}

// New creates a new config, applying the initialization functions,
// filling gaps and validating its content.
func New[T any](options ...func(*T) error) (*T, error) {
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
