package dnsserver

import (
	"crypto/tls"
	"net"
)

// HasSecure tells if the [Server] will handle DoT requests.
func (srv *Server) HasSecure() bool {
	return srv.cfg.TLSConfig != nil
}

// asSecureListeners converts a slice of TCP Listeners into TLS Listeners.
func (srv *Server) asSecureListeners(listeners []*net.TCPListener) []net.Listener {
	var out []net.Listener

	if l := len(listeners); l > 0 {
		tlsConf := srv.cfg.TLSConfig

		out = make([]net.Listener, l)
		for i, tcp := range listeners {
			out[i] = tls.NewListener(tcp, tlsConf)
		}
	}
	return out
}
