package sidecar

import (
	"darvaza.org/core"
	"darvaza.org/x/net/bind"
)

func (srv *Server) initAddresses() error {
	// convert interfaces to addresses
	da := &srv.cfg.Addresses
	if len(da.Interfaces) > 0 {
		s, err := core.GetIPAddresses(da.Interfaces...)
		switch {
		case len(s) > 0:
			da.Addresses = append(da.Addresses, s...)
		case err != nil:
			return err
		}

		da.Interfaces = []string{}
	}
	return nil
}

// Listen listens to all needed ports
func (srv *Server) Listen() error {
	keepalive := srv.cfg.Addresses.KeepAlive
	lc := bind.NewListenConfig(srv.eg.Context(), keepalive)
	return srv.ListenWithListener(lc)
}

// ListenWithUpgrader listens to all needed ports using a ListenUpgrader
// like tableflip
func (srv *Server) ListenWithUpgrader(upg bind.Upgrader) error {
	keepalive := srv.cfg.Addresses.KeepAlive
	lc := bind.NewListenConfig(srv.eg.Context(), keepalive)
	return srv.ListenWithListener(lc.WithUpgrader(upg))
}

// ListenWithListener listens to all needed ports using a net.ListenerConfig
// context
func (srv *Server) ListenWithListener(lc bind.TCPUDPListener) error {
	err := srv.hs.ListenWithListener(lc)
	if err != nil {
		return err
	}

	err = srv.ds.ListenWithListener(lc)
	if err != nil {
		_ = srv.hs.Close()
		return err
	}

	return nil
}
