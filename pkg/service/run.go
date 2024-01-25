package service

import (
	"context"
	"time"

	"darvaza.org/core"
)

// doStart starts the service worker when started
// by the service manager
func (s *Service) doStart() error {
	if err := s.spawnLogHandler(); err != nil {
		return err
	}

	s.spawnService()

	if d := s.Config.SanityDelay; d > 0 {
		select {
		case <-time.After(d):
			// done waiting
			return nil
		case <-s.Cancelled():
			// failed while waiting, wait until it finished
			// shutting down.
			return s.Wait()
		}
	}

	return s.Err()
}

func (s *Service) spawnService() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			if err := core.AsRecovered(recover()); err != nil {
				s.doCancel(err)
			}
		}()

		err := s.run(s.serve, s.args)
		s.doCancel(err)
	}()
}

// doStop stops the service worker when started
// by the service manager
func (s *Service) doStop() error {
	s.doCancel(nil)
	return s.Wait()
}

func (s *Service) doCancel(cause error) {
	if cause == nil {
		cause = context.Canceled
	}

	if s.cancelled.CompareAndSwap(nil, cause) {
		s.cancel()
	}
}

// Wait waits until the worker has finished
// and returns its error value
func (s *Service) Wait() error {
	s.wg.Wait()
	return s.Err()
}

// Cancelled returns a channel that is closed when
// shutdown has been initiated.
func (s *Service) Cancelled() <-chan struct{} {
	return s.ctx.Done()
}

// Done returns a channel that is closed when the
// worker has finished.
func (s *Service) Done() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		s.wg.Wait()
	}()
	return ch
}

// Err returns the recorded cancellation cause.
func (s *Service) Err() error {
	if err, ok := s.cancelled.Load().(error); ok {
		return err
	}

	return nil
}
