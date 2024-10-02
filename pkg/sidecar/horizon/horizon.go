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
	"github.com/miekg/dns"
)

// Match specifies how the Remote made it through,
// and it's stored in the request's context.
type Match struct {
	Horizon    string
	RemoteAddr netip.Addr
	CIDR       netip.Prefix
}

// IsValid checks if the [Match] contains consistent information
func (m Match) IsValid() bool {
	if m.RemoteAddr.IsValid() {
		return m.CIDR.Contains(m.RemoteAddr)
	}
	return false
}

// Horizons is a list of all known horizons sorted by
// priority.
type Horizons struct {
	s []*Horizon
	n map[string]*Horizon

	ExchangeContextFunc func(netip.Addr, *dns.Msg) context.Context
	ExchangeContext     context.Context
	ExchangeTimeoutFunc func(netip.Addr, *dns.Msg) time.Duration
	ExchangeTimeout     time.Duration

	ContextKey *core.ContextKey[Match]
}

// AppendNew creates a [Horizon] based on a [Config] and endpoints.
// [Config.Name] must be unique.
func (s *Horizons) AppendNew(hc Config, h http.Handler, e resolver.Exchanger) error {
	z := hc.New(h, e)
	return s.Append(z)
}

// Append attaches an existing [Horizon]. Name must be unique.
func (s *Horizons) Append(z *Horizon) error {
	if z == nil {
		return core.ErrInvalid
	}

	if s.n == nil {
		s.n = make(map[string]*Horizon)
	}

	if _, ok := s.n[z.n]; ok {
		return core.Wrap(core.ErrExists, "%s", z.String())
	}

	s.s = append(s.s, z)
	s.n[z.n] = z
	return nil
}

// Len returns the number of defined horizons
func (s *Horizons) Len() int {
	return len(s.s)
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

// Get finds a Horizon by name.
func (s Horizons) Get(name string) *Horizon {
	if s.n != nil {
		if z, ok := s.n[name]; ok {
			return z
		}
	}

	return nil
}

// Horizon is one horizon
type Horizon struct {
	n string
	r []netip.Prefix

	h http.Handler
	e resolver.Exchanger
}

// Name returns the name of the [Horizon]
func (z *Horizon) Name() string {
	return z.n
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
