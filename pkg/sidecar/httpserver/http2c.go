package httpserver

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"darvaza.org/middleware"
)

// HasInsecure tells if the [Server] will handle plain HTTP
// requests.
func (srv *Server) HasInsecure() bool {
	return srv.cfg.Bind.AllowInsecure || srv.cfg.Bind.EnableInsecure
}

// NewH2CServer creates a new H2C capable [http.Server].
func (srv *Server) NewH2CServer(h http.Handler, addr net.Addr) *http.Server {
	h1s := srv.NewHTTPServer("h2c", addr)
	h2s := &http2.Server{
		IdleTimeout: srv.cfg.IdleTimeout,
	}

	h1s.Handler = h2c.NewHandler(h, h2s)
	return h1s
}

// NewH2CHandler returns the [http.Handler] to use on the H2C server.
func (srv *Server) NewH2CHandler(h http.Handler) http.Handler {
	switch {
	case !srv.cfg.Bind.AllowInsecure:
		// only ACME-HTTP-01 and https redirect
		h = srv.NewHTTPSRedirectHandler()
	case h == nil:
		// no handler implies 404.
		h = http.NotFoundHandler()
	}

	// ACME-HTTP-01 handler or 404 for /.well-known/acme-challenge
	h = AcmeHTTP01Middleware(h, srv.cfg.AcmeHTTP01)

	// Advertise QUIC
	h = srv.QUICHeadersMiddleware(h)

	return h
}

// NewHTTPSRedirectHandler creates a new handler that redirects everything to
// https.
func (srv *Server) NewHTTPSRedirectHandler() http.Handler {
	port := srv.cfg.Bind.Port
	h := middleware.NewHTTPSRedirectHandler(int(port))

	return h
}

func (srv *Server) spawnH2C(h http.Handler, listeners []*net.TCPListener, graceful time.Duration) error {
	// wrap
	h = srv.NewH2CHandler(h)

	for _, lsn := range listeners {
		s := srv.NewH2CServer(h, lsn.Addr())

		srv.spawnTCP(s, "h2c", lsn, graceful)
	}

	return nil
}
