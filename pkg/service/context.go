package service

import (
	"context"
	"errors"

	"darvaza.org/core"
)

// NewContext prepares a context for cmd.ExecuteContext() or cmd.SetContext()
func (s *Service) NewContext(ctx context.Context) (context.Context, error) {
	switch {
	case s == nil, ctx == nil:
		// invalid call
		return nil, core.ErrInvalid
	case s.ctx != nil:
		// service already initialized
		return nil, core.ErrExists
	default:
		// prepare context to be used with this service
		s2, ok := svcCtxKey.Get(ctx)
		switch {
		case !ok:
			// store ourselves
			ctx = svcCtxKey.WithValue(ctx, s)
			return ctx, nil
		case s2 != s:
			// conflict
			err := errors.New("context already contains a different Service")
			return nil, err
		default:
			// already stored
			return ctx, nil
		}
	}
}

func (s *Service) initContext(ctx context.Context) error {
	// context
	ctx2, err := s.NewContext(ctx)
	switch {
	case err != nil:
		// bad call
		return err
	case ctx != ctx2:
		// command started with a uninitialized context
		return core.ErrInvalid
	default:
		// bind
		s.ctx, s.cancel = context.WithCancel(ctx)
		return nil
	}
}

// GetService gets the current [Service] from the context.
func GetService(ctx context.Context) (*Service, bool) {
	return svcCtxKey.Get(ctx)
}

// WithService stores the given [Service] on a sub-context.
func WithService(ctx context.Context, svc *Service) context.Context {
	return svcCtxKey.WithValue(ctx, svc)
}

var svcCtxKey = core.NewContextKey[*Service]("service")
