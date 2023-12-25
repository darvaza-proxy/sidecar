package horizon

import (
	"net/http"
	"net/netip"

	"darvaza.org/core"
	"darvaza.org/resolver"
)

// Config describe a Horizon
type Config struct {
	Name   string
	Ranges []netip.Prefix

	Middleware         func(http.Handler) http.Handler
	ExchangeMiddleware func(resolver.Exchanger) resolver.Exchanger
}

// New assembles a new [Horizon] using the Config and the given entrypoints
func (hc *Config) New(h http.Handler, e resolver.Exchanger) *Horizon {
	z := &Horizon{
		n: hc.Name,
		r: hc.Ranges,
	}

	if h == nil {
		z.h = http.HandlerFunc(ForbiddenHTTP)
	} else if fn := hc.Middleware; fn != nil {
		z.h = fn(h)
	} else {
		z.h = h
	}

	if e == nil {
		z.e = resolver.ExchangerFunc(ForbiddenExchange)
	} else if fn := hc.ExchangeMiddleware; fn != nil {
		z.e = fn(e)
	} else {
		z.e = e
	}

	return z
}

// Configs represents a sorted list of Horizon configurations
type Configs []Config

// New assembles a new [Horizons] using the [Configs] list and the given entrypoints
func (hcc Configs) New(h http.Handler, e resolver.Exchanger) (*Horizons, error) {
	s := new(Horizons)
	for _, hc := range hcc {
		err := s.AppendNew(hc, h, e)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Must assembles a new [Horizons] using the [Configs] list but panics
// if there is an error.
func (hcc Configs) Must(h http.Handler, e resolver.Exchanger) *Horizons {
	s, err := hcc.New(h, e)
	if err != nil {
		core.PanicWrap(err, "horizon.Configs")
	}

	return s
}
