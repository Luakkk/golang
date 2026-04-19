package service

// ADVANCED bonus: Context-aware versions of UserService methods.
// Each method checks ctx.Done() before (and concurrently with) the repo call,
// so cancellation and deadline-exceeded signals are respected.

import (
	"context"

	"practice-8/repository"
)

// GetUserByIDCtx is a context-aware version of GetUserByID.
// It returns immediately with ctx.Err() if the context is already done,
// and also stops waiting if the context is cancelled while the repo call runs.
func (s *UserService) GetUserByIDCtx(ctx context.Context, id int) (*repository.User, error) {
	// Fast path: context already cancelled/timed-out before we even start.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	type result struct {
		user *repository.User
		err  error
	}

	ch := make(chan result, 1)
	go func() {
		user, err := s.repo.GetUserByID(id)
		ch <- result{user, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-ch:
		return res.user, res.err
	}
}
