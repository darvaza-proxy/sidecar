package service

import (
	"time"
)

func (s *Service) doStart() error {
	var (
		ctx  = s.ctx
		cmd  = &s.cmd
		args = s.cmd.Flags().Args()
		done = make(chan struct{})
	)

	if err := s.cmd.ValidateArgs(args); err != nil {
		// invalid arguments
		return err
	}

	s.preRun()

	if err := s.prepare(ctx, cmd, args); err != nil {
		// early failure
		return err
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer s.postRun()
		defer close(done)

		err := s.run(ctx, cmd, args)
		s.err.Store(err)
	}()

	select {
	case <-done:
		// early run error
	case <-time.After(s.d):
		// enough wait, it is running
	}

	return s.Err()
}

func (s *Service) doStop() error {
	s.cancelOnce.Do(s.cancel)
	return s.Wait()
}

// Wait waits until the worker has finished
// and returns its error value
func (s *Service) Wait() error {
	s.wg.Wait()
	return s.Err()
}

// Err returns the recorded run error
func (s *Service) Err() error {
	if err, ok := s.err.Load().(error); ok {
		return err
	}

	return nil
}
