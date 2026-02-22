package usecase

import (
	"context"

	"practice3-go/internal/repository"
	"practice3-go/pkg/modules"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUsers(ctx context.Context) ([]modules.User, error) {
	return u.repo.GetUsers(ctx)
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id int) (*modules.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u *UserUsecase) CreateUser(ctx context.Context, in modules.CreateUserInput) (int, error) {
	return u.repo.CreateUser(ctx, in)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id int, in modules.UpdateUserInput) error {
	return u.repo.UpdateUser(ctx, id, in)
}

func (u *UserUsecase) DeleteUserByID(ctx context.Context, id int) (int64, error) {
	return u.repo.DeleteUserByID(ctx, id)
}