package dnsserver

import (
	"fmt"
	"net/netip"

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

	if log, ok := srv.info().WithEnabled(); ok {
		log.WithFields(slog.Fields{
			"LocalAddr": ap.String(),
			"Proto":     fmt.Sprintf("dns:%s", proto),
		}).Printf("Listening %s", genListening(proto, ap))
	}
}

func genListening(proto string, ap netip.AddrPort) string {
	var u string

	switch proto {
	case "tcp", "udp":
		var host string

		if ap.Port() != DefaultInsecurePort {
			host = ap.String()
		} else {
			host = ap.Addr().String()
		}

		u = fmt.Sprintf("%s@%s (%s)", "", host, proto)
	case "tcp+tls":
		var host string

		if ap.Port() != DefaultSecurePort {
			host = ap.String()
		} else {
			host = ap.Addr().String()
		}

		u = fmt.Sprintf("%s@%s (%s)", "+tls ", host, "tcp")
	default:
		u = fmt.Sprintf("%s://%s", proto, ap.String())
	}

	return u
}

func (srv *Server) logShuttingDown(proto string, ap netip.AddrPort) {
	if log, ok := srv.debug().WithEnabled(); ok {
		log.WithFields(slog.Fields{
			"LocalAddr": ap.String(),
			"Proto":     fmt.Sprintf("dns:%s", proto),
		}).Print("Shutting down")
	}
}
