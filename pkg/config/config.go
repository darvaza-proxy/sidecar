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
func LoadFile(filename string, v any) error {
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

	err = dec.Decode(data, v)
	if err != nil {
		// failed to decode
		err = &os.PathError{
			Path: filename,
			Op:   "decode",
			Err:  err,
		}
		return err
	}

	err = config.Prepare(v)
	if err != nil {
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

// Prepare fills any gap in the object and validates its content.
func Prepare(v any) error {
	return config.Prepare(v)
}
