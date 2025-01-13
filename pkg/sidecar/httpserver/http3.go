package httpserver

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"darvaza.org/core"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

const (
	// AltSvcHeader is the header label used to advertise
	// QUIC support
	AltSvcHeader = "Alt-Svc"

	// GrabQUICHeadersRetry indicates how long we wait to
	// grab the generated Alt-Svc header
	GrabQUICHeadersRetry = 10 * time.Millisecond
)

// asQUICEarlyListeners converts a slice of UDP Listeners into QUIC Listeners.
func (srv *Server) asQUICEarlyListeners(listeners []*net.UDPConn) ([]*quic.EarlyListener, error) {
	var out []*quic.EarlyListener

	if l := len(listeners); l > 0 {
		cfg := srv.NewQUICConfig()
		tlsConf := srv.NewTLSConfig()
		tlsConf = http3.ConfigureTLSConfig(tlsConf)

		out = make([]*quic.EarlyListener, l)
		for i, udp := range listeners {
			lsn, err := quic.ListenEarly(udp, tlsConf, cfg)
			if err != nil {
				return nil, err
			}
			out[i] = lsn
		}
	}

	return out, nil
}

// NewH3Server creates a new [http3.Server].
func (*Server) NewH3Server(h http.Handler, addr net.Addr) *http3.Server {
	if h == nil {
		h = http.NotFoundHandler()
	}

	return &http3.Server{
		Addr:    addr.String(),
		Handler: h,
	}
}

// NewH3Handler returns the [http.Handler] to use on the H3 server.
func (srv *Server) NewH3Handler(h http.Handler) http.Handler {
	if h == nil {
		h = http.NotFoundHandler()
	}

	// ACME-HTTP-01 handler or 404 for /.well-known/acme-challenge
	h = AcmeHTTP01Middleware(h, srv.cfg.AcmeHTTP01)

	return h
}

func (srv *Server) spawnH3(h http.Handler, listeners []*quic.EarlyListener, graceful time.Duration) error {
	h = srv.NewH3Handler(h)

	for _, lsn := range listeners {
		h3s := srv.NewH3Server(h, lsn.Addr())

		srv.spawnQUIC(h3s, lsn, graceful)
	}

	return nil
}

func (srv *Server) spawnQUIC(h3s *http3.Server, lsn *quic.EarlyListener, _ time.Duration) {
	const proto = "h3"

	addr, ok := core.AddrPort(lsn.Addr())
	if !ok {
		core.Panic("unreachable")
	}

	srv.eg.Go(func(_ context.Context) error {
		srv.logListening(proto, addr)
		return h3s.ServeListener(lsn)
	}, func() error {
		srv.logShuttingDown(proto, addr)
		// err := h3s.CloseGracefully(graceful) // not implemented
		return h3s.Close()
	})

	srv.eg.Go(func(ctx context.Context) error {
		return srv.grabQUICHeaders(ctx, h3s)
	}, nil)
}

// NewQUICConfig returns the [quic.Config] to be used on the
// [http3.Server].
func (*Server) NewQUICConfig() *quic.Config {
	return &quic.Config{}
}

// SetQUICHeaders appends QUIC's Alt-Svc to the [http.Response] headers.
func (srv *Server) SetQUICHeaders(hdr http.Header) error {
	if s := srv.getQUICAltSvc(); s != "" {
		hdr[AltSvcHeader] = append(hdr[AltSvcHeader], s)
	}
	return http3.ErrNoAltSvcPort
}

// QUICHeadersMiddleware creates a middleware function
// that injects Alt-Svc on the [http.Response] headers.
func (srv *Server) QUICHeadersMiddleware(next http.Handler) http.Handler {
	h := func(rw http.ResponseWriter, req *http.Request) {
		_ = srv.SetQUICHeaders(rw.Header())
		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(h)
}

// grabQUICHeader tries periodically to grab the Alt-Svc headers corresponding
// to a server until it succeeds or the given context is cancelled.
func (srv *Server) grabQUICHeaders(ctx context.Context, h3s *http3.Server) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(GrabQUICHeadersRetry):
			hdr := make(http.Header)

			if err := h3s.SetQUICHeaders(hdr); err == nil {
				// success
				srv.appendQUICHeaders(hdr[AltSvcHeader])
				return nil
			}
		}
	}
}

func (srv *Server) getQUICAltSvc() string {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	return srv.quicAltSvc
}

func (srv *Server) appendQUICHeaders(alts []string) {
	var s []string

	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.quicAltSvc != "" {
		s = strings.Split(srv.quicAltSvc, ",")
	}

	for i, hdr := range alts {
		srv.debug().Printf("%s[%v]: %s", AltSvcHeader, i, hdr)

		for _, part := range strings.Split(hdr, ",") {
			part = strings.TrimSpace(part)

			if !core.SliceContains(s, part) {
				s = append(s, part)
			}
		}
	}

	srv.quicAltSvc = strings.Join(s, ",")
}
