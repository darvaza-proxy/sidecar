package config

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v3"
)

// YAML is a Encoder/Decoder for YAML
type YAML struct {
	Spaces int
}

// Decode takes YAML encoded data to populate a data structure
func (*YAML) Decode(data string, v any) error {
	return yaml.Unmarshal([]byte(data), v)
}

// SetIndent specifies what to use for indenting the serialized
// YAML content. Tabs are expanded to 8 spaces.
func (p *YAML) SetIndent(s string) bool {
	var spaces int

	for _, c := range s {
		if c == '\t' {
			spaces += 8
		} else {
			spaces++
		}
	}

	p.Spaces = spaces
	return true
}

// WriteTo encodes data into a given [io.Writer]. If encoding fails,
// the writer won't be touched.
func (p *YAML) WriteTo(v any, w io.Writer) (int64, error) {
	var buf bytes.Buffer

	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(p.Spaces)
	if err := enc.Encode(v); err != nil {
		return 0, err
	}

	return buf.WriteTo(w)
}

// NewYAML creates a new [YAML] Encoder/Decoder with
// two spaces as default indentation.
func NewYAML() *YAML {
	return &YAML{
		Spaces: 2,
	}
}

func init() {
	register("yaml",
		func() Decoder { return NewYAML() },
		func() Encoder { return NewYAML() },
		"yml")
}
