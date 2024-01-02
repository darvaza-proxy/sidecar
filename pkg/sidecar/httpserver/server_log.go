package httpserver

import (
	"fmt"
	"net/netip"
	"net/url"

	"darvaza.org/core"
	"darvaza.org/slog"
)

func (srv *Server) info() slog.Logger {
	return srv.cfg.Logger.Info()
}

func (srv *Server) debug() slog.Logger {
	return srv.cfg.Logger.Debug()
}

func (srv *Server) logListening(proto string, ap netip.AddrPort) {
	addr := ap.Addr()

	if addr.IsUnspecified() {
		port := ap.Port()

		// all
		addrs, _ := core.GetIPAddresses()
		if len(addrs) > 0 {
			for _, ip := range addrs {
				ap := netip.AddrPortFrom(ip, port)
				srv.logListening(proto, ap)
			}
			return
		}
	}

	if l, ok := srv.info().WithEnabled(); ok {
		l.WithFields(slog.Fields{
			"LocalAddr": ap.String(),
			"Proto":     proto,
		}).Printf("Listening %s", genListening(proto, ap))
	}
}

func genListening(proto string, ap netip.AddrPort) string {
	var defaultPort uint16
	var u url.URL
	var s string

	addr := ap.Addr()
	port := ap.Port()

	switch proto {
	case "h2c":
		defaultPort = DefaultInsecurePort
		u.Scheme = "http"
	default:
		defaultPort = DefaultSecurePort
		u.Scheme = "https"
	}

	if port == defaultPort {
		u.Host = addr.String()
	} else {
		u.Host = ap.String()
	}

	s = u.String()

	if proto == "h3" {
		s = fmt.Sprintf("%s (%s)", s, "udp")
	}

	return s
}

func (srv *Server) logShuttingDown(proto string, ap netip.AddrPort) {
	if l, ok := srv.debug().WithEnabled(); ok {
		l.WithFields(slog.Fields{
			"LocalAddr": ap.String(),
			"Proto":     proto,
		}).Print("Shutting down")
	}
}
