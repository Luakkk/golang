package service

// Task 2 + BONUS EASY: Unit tests for UserService using gomock mocks.
// All test functions use table-driven tests with t.Run() subtests,
// testify/assert for soft assertions, and testify/require for hard stops.

import (
	"errors"
	"testing"

	"practice-8/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// ─── From tutorial Step 3 ────────────────────────────────────────────────────

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.CreateUser(user)
	assert.NoError(t, err)
}

// ─── Task 2 + BONUS EASY: RegisterUser (table-driven) ────────────────────────

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		setup     func(m *repository.MockUserRepository)
		wantErrContains string
	}{
		{
			name:  "user already exists",
			email: "existing@mail.com",
			setup: func(m *repository.MockUserRepository) {
				m.EXPECT().GetByEmail("existing@mail.com").
					Return(&repository.User{ID: 1, Name: "Existing"}, nil)
			},
			wantErrContains: "user with this email already exists",
		},
		{
			name:  "new user success",
			email: "new@mail.com",
			setup: func(m *repository.MockUserRepository) {
				m.EXPECT().GetByEmail("new@mail.com").Return(nil, nil)
				m.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
			wantErrContains: "",
		},
		{
			name:  "repository error on CreateUser",
			email: "new@mail.com",
			setup: func(m *repository.MockUserRepository) {
				m.EXPECT().GetByEmail("new@mail.com").Return(nil, nil)
				m.EXPECT().CreateUser(gomock.Any()).Return(errors.New("db connection failed"))
			},
			wantErrContains: "db connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockUserRepository(ctrl)
			tt.setup(mockRepo)
			svc := NewUserService(mockRepo)

			err := svc.RegisterUser(&repository.User{ID: 2, Name: "New"}, tt.email)

			if tt.wantErrContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ─── Task 2 + BONUS EASY: UpdateUserName (table-driven) ──────────────────────

func TestUpdateUserName(t *testing.T) {
	tests := []struct {
		name            string
		userID          int
		newName         string
		setup           func(m *repository.MockUserRepository)
		wantErrContains string
		// For case 4: verify the user.Name was changed BEFORE UpdateUser is called
		wantNameChanged bool
	}{
		{
			name:            "empty name",
			userID:          2,
			newName:         "",
			setup:           func(m *repository.MockUserRepository) {},
			wantErrContains: "name cannot be empty",
		},
		{
			name:    "user not found / repo error",
			userID:  99,
			newName: "NewName",
			setup: func(m *repository.MockUserRepository) {
				m.EXPECT().GetUserByID(99).Return(nil, errors.New("user not found"))
			},
			wantErrContains: "user not found",
		},
		{
			name:    "successful update",
			userID:  2,
			newName: "UpdatedName",
			setup: func(m *repository.MockUserRepository) {
				m.EXPECT().GetUserByID(2).Return(&repository.User{ID: 2, Name: "OldName"}, nil)
				m.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
			wantErrContains: "",
		},
		{
			// Verify that user.Name is actually changed before UpdateUser is called.
			// Passing the expected *User directly — gomock uses deep equality, so this
			// only matches if Name == "ChangedName" at the moment of the call.
			name:            "UpdateUser fails — verify name was changed first",
			userID:          2,
			newName:         "ChangedName",
			wantNameChanged: true,
			setup: func(m *repository.MockUserRepository) {
				m.EXPECT().GetUserByID(2).Return(&repository.User{ID: 2, Name: "OldName"}, nil)
				m.EXPECT().UpdateUser(&repository.User{ID: 2, Name: "ChangedName"}).
					Return(errors.New("update failed"))
			},
			wantErrContains: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockUserRepository(ctrl)
			tt.setup(mockRepo)
			svc := NewUserService(mockRepo)

			err := svc.UpdateUserName(tt.userID, tt.newName)

			if tt.wantErrContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ─── Task 2 + BONUS EASY: DeleteUser (table-driven) ──────────────────────────

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name            string
		userID          int
		setup           func(m *repository.MockUserRepository)
		wantErrContains string
	}{
		{
			name:   "attempt to delete admin user (id=1)",
			userID: 1,
			setup:  func(m *repository.MockUserRepository) {
				// DeleteUser on repo must NOT be called
			},
			wantErrContains: "it is not allowed to delete admin user",
		},
		{
			name:   "successful delete",
			userID: 5,
			setup: func(m *repository.MockUserRepository) {
				// Verify that DeleteUser IS called with the correct id
				m.EXPECT().DeleteUser(5).Return(nil)
			},
			wantErrContains: "",
		},
		{
			name:   "repository error",
			userID: 7,
			setup: func(m *repository.MockUserRepository) {
				// Verify that DeleteUser IS called even when it returns an error
				m.EXPECT().DeleteUser(7).Return(errors.New("db error"))
			},
			wantErrContains: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockUserRepository(ctrl)
			tt.setup(mockRepo)
			svc := NewUserService(mockRepo)

			err := svc.DeleteUser(tt.userID)

			if tt.wantErrContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
