// Package horizon implements entrypoints based on the network
// they belong to
package horizon

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"
	"time"

	"darvaza.org/core"
	"darvaza.org/resolver"
)

// Match specifies how the Remote made it through,
// and it's stored in the request's context.
type Match struct {
	Horizon    string
	RemoteAddr netip.Addr
	CIDR       netip.Prefix
}

// Horizons is a list of all known horizons sorted by
// priority.
type Horizons struct {
	s []*Horizon

	ExchangeContext context.Context
	ExchangeTimeout time.Duration

	ContextKey *core.ContextKey[Match]
}

// AppendNew creates a [Horizon] based on a [Config] and endpoints
func (s *Horizons) AppendNew(hc Config, h http.Handler, e resolver.Exchanger) {
	s.s = append(s.s, hc.New(h, e))
}

// Append attaches an existing [Horizon]
func (s *Horizons) Append(z *Horizon) {
	s.s = append(s.s, z)
}

// Match finds the CIDR and Horizon corresponding to the given address
func (s Horizons) Match(addr netip.Addr) (*Horizon, netip.Prefix, bool) {
	for i := range s.s {
		z := s.s[i]

		if cidr, ok := z.Match(addr); ok {
			return z, cidr, true
		}
	}

	return nil, netip.Prefix{}, false
}

// Horizon is one horizon
type Horizon struct {
	n string
	r []netip.Prefix

	h http.Handler
	e resolver.Exchanger
}

func (z *Horizon) String() string {
	return fmt.Sprintf("%s:%q", "Horizon", z.n)
}

// SetDefaults fills gaps in the [Horizon]
func (z *Horizon) SetDefaults() error {
	if z.r == nil {
		z.r = []netip.Prefix{}
	}

	if z.h == nil {
		z.h = http.HandlerFunc(ForbiddenHTTP)
	}

	if z.e == nil {
		z.e = resolver.ExchangerFunc(ForbiddenExchange)
	}

	return nil
}

// Match finds the first matching CIDR for the given address
func (z *Horizon) Match(addr netip.Addr) (netip.Prefix, bool) {
	if len(z.r) == 0 {
		// any
		if addr.Is6() {
			addr = netip.IPv6Unspecified()
		} else {
			addr = netip.IPv4Unspecified()
		}

		r := netip.PrefixFrom(addr, 0)
		return r, true
	}

	for _, r := range z.r {
		if r.Contains(addr) {
			// match
			return r, true
		}
	}

	// none
	return netip.Prefix{}, false
}

// InRange checks if the remote address belongs on this [Horizon]
func (z *Horizon) InRange(addr netip.Addr) bool {
	_, ok := z.Match(addr)
	return ok
}