package config

// Autodetect is a Decoder that tries all Decoders
type Autodetect struct{}

// Decode takes encoded data and tries all registered decoders.
// Decode returns [ErrUnknownFormat] if all decoders fail.
func (p *Autodetect) Decode(data string, v any) error {
	for name, r := range registry {
		if name != "auto" {
			if p.tryDecode(r, data, v) {
				// success
				return nil
			}
		}
	}

	return ErrUnknownFormat
}

func (*Autodetect) tryDecode(r *registryEntry, data string, v any) bool {
	f := r.NewDecoder
	if f == nil {
		return false
	}

	dec := f()
	if dec == nil {
		return false
	}

	err := dec.Decode(data, v)
	return err == nil
}

func init() {
	register("auto", func() Decoder { return &Autodetect{} }, nil)
}
