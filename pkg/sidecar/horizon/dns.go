package horizon

import (
	"context"
	"net/netip"
	"time"

	"github.com/miekg/dns"

	"darvaza.org/core"
	"darvaza.org/resolver"
	"darvaza.org/resolver/pkg/errors"
)

var (
	_ dns.Handler        = (*Horizons)(nil)
	_ resolver.Exchanger = (*Horizons)(nil)
	_ resolver.Exchanger = (*Horizon)(nil)
)

// ServeDNS implements the [dns.Handler] interface
func (s Horizons) ServeDNS(rw dns.ResponseWriter, req *dns.Msg) {
	z, m, ok := s.MatchDNSRequest(rw)
	if !ok {
		HandleForbiddenExchange(rw, req)
		return
	}

	ctx, cancel := s.newDNSLookupContext(m, req)
	defer cancel()

	rsp, err := z.Exchange(ctx, req)
	if err != nil {
		rsp = errors.ErrorAsMsg(req, err)
	}

	_ = rw.WriteMsg(rsp)
}

func (s Horizons) newDNSLookupContext(m Match, req *dns.Msg) (context.Context, context.CancelFunc) {
	var ctx context.Context
	var timeout time.Duration

	// parent
	if s.ExchangeContextFunc != nil {
		ctx = s.ExchangeContextFunc(m.RemoteAddr, req)
	} else {
		ctx = s.ExchangeContext
	}

	if ctx == nil {
		ctx = context.Background()
	}

	// attach match
	ctx = s.ContextKey.WithValue(ctx, m)

	// timeout
	if s.ExchangeTimeoutFunc != nil {
		timeout = s.ExchangeTimeoutFunc(m.RemoteAddr, req)
	} else {
		timeout = s.ExchangeTimeout
	}

	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return ctx, func() {}
}

// Exchange implements the [resolver.Exchanger] interface but requires a [Match]
// in the context
func (s Horizons) Exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	if s.ContextKey != nil {
		m, ok := s.ContextKey.Get(ctx)
		if ok && m.IsValid() {
			if z := s.Get(m.Horizon); z != nil {
				return z.Exchange(ctx, req)
			}
		}
	}
	return ForbiddenExchange(ctx, req)
}

// MatchDNSRequest find the Horizon corresponding to an [http.Request] and
// prepares a [Match] to include in the context.
func (s Horizons) MatchDNSRequest(rw dns.ResponseWriter) (*Horizon, Match, bool) {
	addr, _ := DNSRemoteAddr(rw)
	if addr.IsValid() {
		z, cidr, ok := s.Match(addr)
		if ok {
			m := Match{
				Horizon:    z.n,
				CIDR:       cidr,
				RemoteAddr: addr,
			}
			return z, m, true
		}
	}

	return nil, Match{}, false
}

// Exchange implements the [resolver.Exchanger] interface
func (z *Horizon) Exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	return z.e.Exchange(ctx, req)
}

// Lookup implements the [resolver.Lookuper] interface
func (z *Horizon) Lookup(ctx context.Context, qName string, qType uint16) (*dns.Msg, error) {
	l := resolver.ExchangerFunc(z.e.Exchange)
	return l.Lookup(ctx, qName, qType)
}

// DNSRemoteAddr extracts the remote address associated with an [dns.ResponseWriter]
func DNSRemoteAddr(rw dns.ResponseWriter) (netip.Addr, error) {
	ap, ok := core.AddrPort(rw.RemoteAddr())
	if !ok || !ap.IsValid() {
		return netip.Addr{}, core.ErrInvalid
	}

	addr := ap.Addr()
	if addr.Is4In6() {
		addr = addr.Unmap()
	}
	return addr, nil
}

// ForbiddenExchange is a resolver.ExchangerFunc that refuses all requests
func ForbiddenExchange(_ context.Context, req *dns.Msg) (*dns.Msg, error) {
	resp := newForbiddenResponse(req)
	return nil, errors.MsgAsError(resp)
}

// HandleForbiddenExchange is a dns.HandlerFunc that refuses all requests
func HandleForbiddenExchange(rw dns.ResponseWriter, req *dns.Msg) {
	resp := newForbiddenResponse(req)
	_ = rw.WriteMsg(resp)
}

func newForbiddenResponse(req *dns.Msg) *dns.Msg {
	resp := new(dns.Msg)
	resp.SetRcode(req, dns.RcodeRefused)
	resp.Compress = false
	resp.RecursionAvailable = true
	return resp
}
