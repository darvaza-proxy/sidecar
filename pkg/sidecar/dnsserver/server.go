// Package dnsserver implements a DNS server for sidecars
package dnsserver

import (
	"context"
	"net"
	"net/netip"
	"syscall"
	"time"

	"github.com/miekg/dns"

	"darvaza.org/core"
)

// Server is a DNS/DoT Server built
// around a shared [core.ErrGroup].
type Server struct {
	cfg Config

	eg  *core.ErrGroup
	sl  *Listeners
	dns []*dns.Server
}

func (ds *Server) setupServer(s *dns.Server) *dns.Server {
	switch {
	case s.TLSConfig != nil:
		s.Net = "tcp+tls"
	case s.Listener != nil:
		s.Net = "tcp"
	default:
		s.Net = "udp"
	}

	s.IdleTimeout = func() time.Duration { return ds.cfg.IdleTimeout }
	s.ReadTimeout = ds.cfg.ReadTimeout
	s.MaxTCPQueries = ds.cfg.MaxTCPQueries

	return s
}

func (ds *Server) prepare(h dns.Handler) error {
	switch {
	case len(ds.dns) > 0:
		// already running
		return core.Wrap(syscall.EBUSY, "server already running")
	case ds.sl == nil:
		// not listening
		return core.Wrap(core.ErrInvalid, "no listeners available")
	}

	if h == nil {
		// NO-OP handler
		h = dns.NewServeMux()
	}

	// 53/TCP
	for _, lsn := range ds.sl.TCP {
		s := ds.setupServer(&dns.Server{
			Listener: lsn,
			Handler:  h,
		})
		ds.dns = append(ds.dns, s)
	}

	// 53/UDP
	for _, lsn := range ds.sl.UDP {
		s := ds.setupServer(&dns.Server{
			PacketConn: lsn,
			Handler:    h,
		})
		ds.dns = append(ds.dns, s)
	}

	// 853/TCP+TLS
	for _, lsn := range ds.sl.TLS {
		s := ds.setupServer(&dns.Server{
			Listener:  lsn,
			TLSConfig: ds.cfg.TLSConfig,
			Handler:   h,
		})
		ds.dns = append(ds.dns, s)
	}

	return nil
}

// Spawn starts all workers and optionally waits a given amount
// to make sure they didn't fail.
func (ds *Server) Spawn(h dns.Handler, wait time.Duration) error {
	if err := ds.prepare(h); err != nil {
		return err
	}

	for _, s := range ds.dns {
		ds.spawnServer(s, ds.cfg.GracefulTimeout)
	}

	if wait > 0 {
		select {
		case <-time.After(wait):
			// done waiting
		case <-ds.eg.Cancelled():
			// failed while waiting
			return ds.eg.Wait()
		}
	}

	return ds.eg.Err()
}

func (ds *Server) spawnServer(s *dns.Server, graceful time.Duration) {
	proto, addr := getServerProtoAddr(s)

	ds.eg.Go(func(_ context.Context) error {
		ds.logListening(proto, addr)
		return s.ActivateAndServe()
	}, func() error {
		ds.logShuttingDown(proto, addr)

		if graceful > 0 {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, graceful)
			defer cancel()

			return s.ShutdownContext(ctx)
		}

		return s.Shutdown()
	})
}

func getServerProtoAddr(s *dns.Server) (string, netip.AddrPort) {
	var addr net.Addr
	var proto = s.Net

	switch {
	case s.PacketConn != nil:
		addr = s.PacketConn.LocalAddr()
	case s.Listener != nil:
		addr = s.Listener.Addr()
	default:
		core.Panic("invalid dns.Server: no listener")
	}

	ap, ok := core.AddrPort(addr)
	if !ok {
		core.Panicf("invalid dns.Server address: %q", addr.String())
	}

	return proto, ap
}

// Serve starts all workers and waits until they have
// finished.
func (ds *Server) Serve(h dns.Handler) error {
	if err := ds.Spawn(h, 0); err != nil {
		// failed to prepare
		return err
	}

	return ds.Wait()
}

// Cancel initiates a cancellation with the given
// reason.
func (ds *Server) Cancel(cause error) {
	ds.eg.Cancel(cause)
}

// Wait waits until all workers have finished.
func (ds *Server) Wait() error {
	return ds.eg.Wait()
}
