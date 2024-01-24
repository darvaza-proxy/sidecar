package httpserver

import (
	"io"
	"net"
	"sync/atomic"
	"syscall"

	"github.com/quic-go/quic-go"

	"darvaza.org/x/net/bind"
)

const (
	// DefaultSecurePort represents the default port for secure HTTP (TCP and UDP)
	DefaultSecurePort = 443
	// DefaultInsecurePort represents the default port for plain HTTP (TCP)
	DefaultInsecurePort = 80
)

// Listeners contains the listeners to be used by this HTTP server.
type Listeners struct {
	Secure   []net.Listener
	Insecure []*net.TCPListener
	Quic     []*quic.EarlyListener

	up atomic.Bool
}

// Close closes all listeners.
func (sl *Listeners) Close() error {
	closeAll(sl.Secure)
	closeAll(sl.Insecure)
	closeAll(sl.Quic)
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

	// HTTPS/QUIC
	if srv.HasSecure() {
		bc.Port = cfg.Port
		bc.DefaultPort = DefaultSecurePort

		tlsLsn, quicLsn, err := srv.bindTLS(bc)
		if err != nil {
			return nil, err
		}
		sl.Secure = tlsLsn
		sl.Quic = quicLsn
	}

	// HTTP
	if srv.HasInsecure() {
		bc.Port = cfg.PortInsecure
		bc.DefaultPort = DefaultInsecurePort
		bc.OnlyTCP = true

		tcp, _, err := bc.Bind()
		if err != nil {
			defer sl.Close()
			return nil, err
		}

		sl.Insecure = tcp
	}

	return sl, nil
}

func (srv *Server) bindTLS(bc *bind.Config) ([]net.Listener, []*quic.EarlyListener, error) {
	tcp, udp, err := bc.Bind()
	if err != nil {
		return nil, nil, err
	}

	tlsLsn := srv.asSecureListeners(tcp)

	quicLsn, err := srv.asQuicEarlyListeners(udp)
	if err != nil {
		closeAll(tlsLsn)
		closeAll(udp)
		return nil, nil, err
	}

	return tlsLsn, quicLsn, nil
}
