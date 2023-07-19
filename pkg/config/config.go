// Package config contains helpers for loading and dumping TOML config files
package config

import (
	"bytes"
	"io"

	"github.com/BurntSushi/toml"

	"darvaza.org/darvaza/shared/config"
	"darvaza.org/darvaza/shared/config/expand"
)

var (
	// Indent is the string used for indenting TOML files when encoding them
	Indent = `  `
)

// LoadFile loads a TOML file by name, expanding environment variables,
// filling gaps and validating its content.
func LoadFile(filename string, v any) error {
	data, err := expand.FromFile(filename, nil)
	if err != nil {
		return err
	}

	_, err = toml.Decode(data, v)
	if err != nil {
		return err
	}

	return config.Prepare(v)
}

// WriteTo writes out the object encoded as TOML.
func WriteTo(v any, w io.Writer) (int64, error) {
	var buf bytes.Buffer

	// encode
	enc := toml.NewEncoder(&buf)
	enc.Indent = Indent
	if err := enc.Encode(v); err != nil {
		return 0, err
	}

	return buf.WriteTo(w)
}

// Prepare fills any gap in the object and validates its content.
func Prepare(v any) error {
	return config.Prepare(v)
}
