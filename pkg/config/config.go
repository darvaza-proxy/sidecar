// Package config contains helpers for loading and dumping TOML config files
package config

import (
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

	dec, ok := NewDecoderByFilename(filename)
	if !ok {
		// fallback to autodetect
		dec, ok = NewDecoder("auto")
	}

	if !ok {
		return ErrUnknownFormat
	}

	err = dec.Decode(data, v)
	if err == nil {
		err = config.Prepare(v)
	}

	return err
}

// Prepare fills any gap in the object and validates its content.
func Prepare(v any) error {
	return config.Prepare(v)
}
