package horizon

import (
	"bytes"
	"net/http"
	"net/netip"

	"darvaza.org/core"
)

var (
	_ http.Handler = (*Horizons)(nil)
	_ http.Handler = (*Horizon)(nil)
)

// ServeHTTP implements the [http.Handler] interface
func (s Horizons) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	z, m, ok := s.MatchHTTPRequest(req)
	if !ok {
		ForbiddenHTTP(rw, req)
		return
	}

	if s.ContextKey != nil {
		// add Match to context
		ctx := req.Context()
		ctx = s.ContextKey.WithValue(ctx, m)
		req = req.WithContext(ctx)
	}

	z.ServeHTTP(rw, req)
}

// MatchHTTPRequest find the Horizon corresponding to an [http.Request] and
// prepares a [Match] to include in the context.
func (s Horizons) MatchHTTPRequest(req *http.Request) (*Horizon, Match, bool) {
	addr, err := HTTPRemoteAddr(req)
	if err == nil {
		z, cidr, ok := s.Match(addr)
		if ok {
			m := Match{
				Horizon:    z.n,
				RemoteAddr: addr,
				CIDR:       cidr,
			}
			return z, m, true
		}
	}

	return nil, Match{}, false
}

// ServeHTTP implements the [http.Handler] interface
func (z *Horizon) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	z.h.ServeHTTP(rw, req)
}

// HTTPRemoteAddr extracts the remote address associated with an [http.Request]
func HTTPRemoteAddr(req *http.Request) (netip.Addr, error) {
	addr, _, err := core.SplitAddrPort(req.RemoteAddr)
	if err != nil {
		return addr, err
	}

	if addr.Is4In6() {
		addr = addr.Unmap()
	}

	return addr, nil
}

// ForbiddenHTTP is an http.HandlerFunc that always return a 403 error
func ForbiddenHTTP(rw http.ResponseWriter, _ *http.Request) {
	msg := bytes.NewBufferString("forbidden")

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusForbidden)
	_, _ = msg.WriteTo(rw)
}
