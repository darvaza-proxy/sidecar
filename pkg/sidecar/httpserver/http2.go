package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"

	"darvaza.org/core"
)

// NewHTTPServer creates a new [http.Server].
func (srv *Server) NewHTTPServer(addr net.Addr) *http.Server {
	return &http.Server{
		Addr: addr.String(),

		ReadTimeout:       srv.cfg.ReadTimeout,
		ReadHeaderTimeout: srv.cfg.ReadHeaderTimeout,
		WriteTimeout:      srv.cfg.WriteTimeout,
		IdleTimeout:       srv.cfg.IdleTimeout,
	}
}

// NewH2Server creates a new HTTP/2 capable [http.Server].
func (srv *Server) NewH2Server(h http.Handler, addr net.Addr) (*http.Server, error) {
	h1s := srv.NewHTTPServer(addr)
	h1s.TLSConfig = srv.NewTLSConfig()
	h1s.Handler = h

	h2s := &http2.Server{}
	if err := http2.ConfigureServer(h1s, h2s); err != nil {
		return nil, err
	}

	return h1s, nil
}

// NewH2Handler returns the [http.Handler] to use on the H2 server.
func (srv *Server) NewH2Handler(h http.Handler) http.Handler {
	if h == nil {
		// no handler implies 404.
		h = http.NotFoundHandler()
	}

	// ACME-HTTP-01 handler or 404 for /.well-known/acme-challenge
	h = AcmeHTTP01Middleware(h, srv.cfg.AcmeHTTP01)

	// Advertise Quic
	h = srv.QuicHeadersMiddleware(h)

	return h
}

func (srv *Server) spawnH2(h http.Handler, listeners []net.Listener, graceful time.Duration) error {
	// wrap
	h = srv.NewH2Handler(h)

	for _, lsn := range listeners {
		s, err := srv.NewH2Server(h, lsn.Addr())
		if err != nil {
			return err
		}

		srv.spawnTCP(s, "h2", lsn, graceful)
	}

	return nil
}

func (srv *Server) spawnTCP(s *http.Server, proto string, lsn net.Listener, graceful time.Duration) {
	addr, ok := core.AddrPort(lsn.Addr())
	if !ok {
		core.Panic("unreachable")
	}

	srv.eg.Go(func(_ context.Context) error {
		srv.logListening(proto, addr)
		return s.Serve(lsn)
	}, func() error {
		srv.logShuttingDown(proto, addr)

		ctx := context.Background()
		if graceful > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, graceful)
			defer cancel()
		}

		return s.Shutdown(ctx)
	})
}
