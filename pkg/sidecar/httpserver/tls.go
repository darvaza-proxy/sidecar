package httpserver

import (
	"crypto/tls"
	"net"
)

// HasSecure tells if the [Server] will handle HTTPS and HTTP/3
// requests.
func (srv *Server) HasSecure() bool {
	return srv.cfg.TLSConfig != nil
}

// NewTLSConfig returns the [tls.Config] to be used on the [Server].
func (srv *Server) NewTLSConfig() *tls.Config {
	if tc := srv.cfg.TLSConfig; tc != nil {
		return tc.Clone()
	}
	return nil
}

// asSecureListeners converts a slice of TCP Listeners into TLS Listeners.
func (srv *Server) asSecureListeners(listeners []*net.TCPListener) []net.Listener {
	var out []net.Listener

	if l := len(listeners); l > 0 {
		tlsConf := srv.NewTLSConfig()

		out = make([]net.Listener, l)
		for i, tcp := range listeners {
			out[i] = tls.NewListener(tcp, tlsConf)
		}
	}
	return out
}
