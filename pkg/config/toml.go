package config

import (
	"bytes"
	"io"

	"github.com/BurntSushi/toml"
)

const (
	// DefaultTOMLIndent is the indentation used unless SetIndent
	// is called on the TOML Encoder
	DefaultTOMLIndent = `  `
)

// TOML is a Encoder/Decoder for TOML
type TOML struct {
	Indent string
}

// Decode takes TOML encoded data to populate a data structure
func (*TOML) Decode(data string, v any) error {
	_, err := toml.Decode(data, v)
	return err
}

// SetIndent specifies what to do for indenting the serialized
// TOML content.
func (p *TOML) SetIndent(indent string) bool {
	p.Indent = indent
	return true
}

// WriteTo encodes data into a given [io.Writer]. If encoding fails,
// the writer won't be touched.
func (p *TOML) WriteTo(v any, w io.Writer) (int64, error) {
	var buf bytes.Buffer

	// encode
	enc := toml.NewEncoder(&buf)
	enc.Indent = p.Indent
	if err := enc.Encode(v); err != nil {
		return 0, err
	}

	return buf.WriteTo(w)
}

// NewTOML creates a new [TOML] Encoder/Decoder with two space
// indentation.
func NewTOML() *TOML {
	return &TOML{
		Indent: DefaultTOMLIndent,
	}
}

func init() {
	register("toml",
		func() Decoder { return NewTOML() },
		func() Encoder { return NewTOML() },
	)
}
