package sidecar

import (
	"darvaza.org/darvaza/shared/net/bind"
)

// Listen listens to all needed ports
func (srv *Server) Listen() error {
	lc := bind.NewListenConfig(srv.ctx, 0)
	return srv.ListenWithListener(lc)
}

// ListenWithUpgrader listens to all needed ports using a ListenUpgrader
// like tableflip
func (srv *Server) ListenWithUpgrader(upg bind.Upgrader) error {
	lc := bind.NewListenConfig(srv.ctx, 0)
	return srv.ListenWithListener(lc.WithUpgrader(upg))
}

// ListenWithListener listens to all needed ports using a net.ListenerConfig
// context
func (srv *Server) ListenWithListener(lc bind.TCPUDPListener) error {
	return srv.hs.ListenWithListener(lc)
}
