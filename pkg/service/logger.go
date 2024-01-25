package service

import (
	"encoding/json"
	"syscall"

	"github.com/kardianos/service"

	"darvaza.org/slog"
	"darvaza.org/slog/handlers/cblog"
)

func (s *Service) spawnLogHandler() error {
	if s.log != nil {
		return syscall.EBUSY
	}

	sl, err := s.ss.SystemLogger(nil)
	if err != nil {
		return err
	}

	ch := make(chan cblog.LogMsg, cblog.DefaultOutputBufferSize)
	s.log, _ = cblog.New(ch)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer close(ch)
		defer func() { _ = recover() }()

		s.runLogHandler(sl, ch)
	}()

	return nil
}

func (s *Service) runLogHandler(sl service.Logger, ch <-chan cblog.LogMsg) {
	for {
		select {
		case m := <-ch:
			s.logEntry(&m, sl)
		case <-s.Cancelled():
			return
		}
	}
}

func (*Service) logEntry(m *cblog.LogMsg, sl service.Logger) {
	// encode
	b, err := json.MarshalIndent(m, ``, `  `)
	msg := string(b)
	switch {
	case err != nil:
		_ = sl.Error("encoder error", err)
	case m.Level < slog.Warn:
		_ = sl.Error(msg)
	case m.Level == slog.Warn:
		_ = sl.Warning(msg)
	default:
		_ = sl.Info(msg)
	}
}

// SystemLogger returns the OS system logger.
// Entries will be recorded JSON encoded.
func (s *Service) SystemLogger() slog.Logger {
	return s.log
}
