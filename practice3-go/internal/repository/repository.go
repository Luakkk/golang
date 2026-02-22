package repository

import (
	"context"

	"practice3-go/internal/repository/_postgres"
	"practice3-go/internal/repository/_postgres/users"
	"practice3-go/pkg/modules"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]modules.User, error)
	GetUserByID(ctx context.Context, id int) (*modules.User, error)
	CreateUser(ctx context.Context, in modules.CreateUserInput) (int, error)
	UpdateUser(ctx context.Context, id int, in modules.UpdateUserInput) error
	DeleteUserByID(ctx context.Context, id int) (int64, error)
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}