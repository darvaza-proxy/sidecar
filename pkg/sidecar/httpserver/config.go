package httpserver

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/netip"
	"time"

	"darvaza.org/core"
	"darvaza.org/darvaza/shared/config"
	"darvaza.org/slog"
	"darvaza.org/slog/handlers/discard"
)

// Config describes how the [Server] will be assembled
// and operate.
type Config struct {
	Context context.Context
	Logger  slog.Logger

	Bind      BindingConfig
	TLSConfig *tls.Config

	// AcmeHTTP01 is an optional [http.Handler] that will
	// receive requests for /.well-known/acme-challenge
	// with a valid token.
	//
	// If no handler is specified, the server will
	// automatically emit a 404 error for requests
	// against this well-known path.
	AcmeHTTP01 http.Handler

	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration

	GracefulTimeout time.Duration
}

// SetDefaults fills gaps in the [Config].
func (sc *Config) SetDefaults() error {
	if sc.Context == nil {
		sc.Context = context.Background()
	}

	if sc.Logger == nil {
		sc.Logger = discard.New()
	}

	return config.Set(sc)
}

// New creates a new [Server] from the [Config], optionally
// taking a shared [core.ErrGroup] for cancellations.
func (sc *Config) New(eg *core.ErrGroup) (*Server, error) {
	if sc == nil {
		sc = new(Config)
	}

	if err := sc.SetDefaults(); err != nil {
		return nil, err
	}

	if eg == nil {
		eg = &core.ErrGroup{
			Parent: sc.Context,
		}
	}

	srv := &Server{
		eg:  eg,
		cfg: *sc,
	}

	return srv, nil
}

// BindingConfig describes what the [Server] will listen.
type BindingConfig struct {
	Addrs        []netip.Addr
	Port         uint16
	PortInsecure uint16

	PortStrict   bool
	PortAttempts int

	// EnableInsecure makes us listen the plain HTTP port
	EnableInsecure bool
	// AllowInsecure makes us handle plain HTTP requests
	// instead of simply redirecting to the HTTPS port.
	AllowInsecure bool
}
