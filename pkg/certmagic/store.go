package certmagic

import (
	"github.com/caddyserver/certmagic"
)

// Store implements a [storage.Store] using
// Caddy's certmagic
type Store struct {
	cmc certmagic.Config
}

// init is run before the options
func (*Store) init() error { return nil }

// prepare is run after the options
func (*Store) prepare() error { return nil }

// New creates a new certmagic [storage.Store] using
// the given options
func New(options ...OptionFunc) (*Store, error) {
	s := &Store{
		cmc: certmagic.Default,
	}

	if err := s.init(); err != nil {
		return nil, err
	}

	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	if err := s.prepare(); err != nil {
		return nil, err
	}

	return s, nil
}
