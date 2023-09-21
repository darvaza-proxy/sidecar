package certmagic

import (
	"context"

	"darvaza.org/core"
	"darvaza.org/slog"
	"darvaza.org/slog/handlers/zap"
)

// OptionFunc ...
type OptionFunc func(*Store) error

// WithContext passes a [context.Context] that can be used to indicate
// cancellation to the [Store].
func WithContext(ctx context.Context) OptionFunc {
	return func(s *Store) error {
		if ctx == nil {
			return core.Wrap(core.ErrInvalid, "no cancel context")
		}

		return core.ErrNotImplemented
	}
}

// WithLogger passes the [slog.Logger] to be used by the [Store].
func WithLogger(log slog.Logger) OptionFunc {
	return func(s *Store) error {
		if log == nil {
			return core.Wrap(core.ErrInvalid, "no logger")
		}

		zl, err := zap.NewReversed(log)
		if err != nil {
			return err
		}
		s.cmc.Logger = zl
		return nil
	}
}

// WithTrustedRoots passes CA certificates the [Store] can trust.
func WithTrustedRoots(rootCerts ...string) OptionFunc {
	return func(s *Store) error {
		ok, err := addTrustedRoots(s, rootCerts)
		switch {
		case err != nil:
			return err
		case !ok:
			return core.Wrap(core.ErrInvalid, "no trusted roots")
		default:
			return nil
		}
	}
}

func addTrustedRoots(s *Store, rootCerts []string) (bool, error) {
	var added bool

	for _, crt := range rootCerts {
		if crt != "" {
			ok, err := s.AddTrustedRoot(crt)
			switch {
			case err != nil:
				return added, err
			case ok:
				added = true
			}
		}
	}

	return added, nil
}

// WithKey passes the x509 key and optionally server certificates
// to be used by the [Store].
func WithKey(key string, certs ...string) OptionFunc {
	return func(s *Store) error {
		if key == "" {
			return core.Wrap(core.ErrInvalid, "no key")
		}

		return s.SetKey(key, certs...)
	}
}

// WithIssuerKey passes an optional x509 key and CA certificate to be used
// when issuing certificates itself
func WithIssuerKey(key string, cert string) OptionFunc {
	return func(s *Store) error {
		if key == "" {
			return core.Wrap(core.ErrInvalid, "no issuer key")
		}

		return s.SetIssuerKey(key, cert)
	}
}
