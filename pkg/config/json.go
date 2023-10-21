package config

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	// DefaultJSONIndent is the indentation used unless SetIndent
	// is called on the JSON Encoder
	DefaultJSONIndent = `  `
)

// JSON is a Encoder/Decoder for JSON
type JSON struct {
	Indent string
}

// Decode takes JSON encoded data to populate a data structure.
func (JSON) Decode(data string, v any) error {
	return json.Unmarshal([]byte(data), v)
}

// SetIndent specifies what to do for indenting the serialized
// JSON content.
func (p *JSON) SetIndent(indent string) bool {
	p.Indent = indent
	return true
}

// WriteTo encodes data into a given [io.Writer]. If encoding fails,
// the writer won't be touched.
func (p *JSON) WriteTo(v any, w io.Writer) (int64, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", p.Indent)

	if err := enc.Encode(v); err != nil {
		return 0, err
	}

	return buf.WriteTo(w)
}

// NewJSON creates a new [JSON] Encoder/Decoder with two space
// indentation.
func NewJSON() *JSON {
	return &JSON{
		Indent: DefaultJSONIndent,
	}
}

func init() {
	register("json",
		func() Decoder { return NewJSON() },
		func() Encoder { return NewJSON() })
}
