package dnsserver

import (
	"io"
	"net"
	"syscall"

	"darvaza.org/darvaza/shared/net/bind"
)

const (
	// DefaultSecurePort is the TCP port used by
	// default to serve DNS over TLS (a.k.a. DoT)
	DefaultSecurePort = 853

	// DefaultInsecurePort is the UDP/TCP port used by default to
	// serve plain DNS
	DefaultInsecurePort = 53
)

// Listeners contains the listeners to be used by this DNS server.
type Listeners struct {
	UDP []*net.UDPConn
	TCP []*net.TCPListener
	TLS []net.Listener
}

// Close closes all listeners.
func (sl *Listeners) Close() error {
	closeAll(sl.UDP)
	closeAll(sl.TCP)
	closeAll(sl.TLS)
	return nil
}

func closeAll[T io.Closer](s []T) {
	for _, l := range s {
		_ = l.Close()
	}
}

// ListenWithListener uses a given [bind.TCPUDPListener] to listen to the
// addresses specified on the [Config].
func (srv *Server) ListenWithListener(lc bind.TCPUDPListener) error {
	if srv.sl != nil {
		return syscall.EBUSY
	}

	cfg := &srv.cfg.Bind

	addrs := make([]string, len(cfg.Addrs))
	for i, addr := range cfg.Addrs {
		addrs[i] = addr.String()
	}

	bc := &bind.Config{
		Addresses:    addrs,
		PortStrict:   cfg.PortStrict,
		PortAttempts: cfg.PortAttempts,
	}
	bc.UseListener(lc)

	sl, err := srv.newListeners(cfg, bc)
	if err != nil {
		return err
	}

	srv.sl = sl
	return nil
}

func (srv *Server) newListeners(cfg *BindingConfig, bc *bind.Config) (*Listeners, error) {
	var sl = new(Listeners)

	// TCP/UDP
	bc.Port = cfg.Port
	bc.Port = DefaultInsecurePort

	tcp, udp, err := bc.Bind()
	if err != nil {
		return nil, err
	}

	sl.TCP = tcp
	sl.UDP = udp

	// DNS over TLS
	if srv.HasSecure() {
		bc.Port = cfg.TLSPort
		bc.DefaultPort = DefaultSecurePort
		bc.OnlyTCP = true

		tcp, _, err = bc.Bind()
		if err != nil {
			defer sl.Close()
			return nil, err
		}

		sl.TLS = srv.asSecureListeners(tcp)
	}

	return sl, nil
}
