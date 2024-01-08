package sidecar

import "darvaza.org/slog"

func (srv *Server) error(err error) slog.Logger {
	l := srv.cfg.Logger.Error()
	if err != nil {
		l = l.WithField(slog.ErrorFieldName, err)
	}
	return l
}

func (srv *Server) warn() slog.Logger {
	return srv.cfg.Logger.Warn()
}
