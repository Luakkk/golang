package service

// ADVANCED bonus: Context support tests for GetUserByIDCtx.
// Tests:
//   1. Cancelled context   — returns context.Canceled immediately, repo NOT called.
//   2. Deadline exceeded   — returns context.DeadlineExceeded, repo NOT called.
//   3. Happy path with ctx — normal repo call succeeds when ctx is active.

import (
	"context"
	"testing"
	"time"

	"practice-8/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserByIDCtx(t *testing.T) {
	tests := []struct {
		name        string
		buildCtx    func() (context.Context, context.CancelFunc)
		setup       func(m *repository.MockUserRepository)
		wantUser    *repository.User
		wantErrIs   error // checked with errors.Is / assert.ErrorIs
	}{
		{
			name: "cancelled context — repo must not be called",
			buildCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // cancel before the call
				return ctx, cancel
			},
			setup:     func(m *repository.MockUserRepository) {}, // no EXPECT — repo not called
			wantUser:  nil,
			wantErrIs: context.Canceled,
		},
		{
			name: "deadline exceeded — repo must not be called",
			buildCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(1 * time.Millisecond) // ensure the deadline has already passed
				return ctx, cancel
			},
			setup:     func(m *repository.MockUserRepository) {}, // no EXPECT — repo not called
			wantUser:  nil,
			wantErrIs: context.DeadlineExceeded,
		},
		{
			name: "active context — successful repo call",
			buildCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			setup: func(m *repository.MockUserRepository) {
				m.EXPECT().
					GetUserByID(42).
					Return(&repository.User{ID: 42, Name: "Alua"}, nil)
			},
			wantUser:  &repository.User{ID: 42, Name: "Alua"},
			wantErrIs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockUserRepository(ctrl)
			tt.setup(mockRepo)
			svc := NewUserService(mockRepo)

			ctx, cancel := tt.buildCtx()
			defer cancel()

			user, err := svc.GetUserByIDCtx(ctx, 42)

			if tt.wantErrIs != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUser, user)
			}
		})
	}
}
