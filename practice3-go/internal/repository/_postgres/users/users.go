package users

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"practice3-go/internal/domain/apperr"
	"practice3-go/internal/repository/_postgres"
	"practice3-go/pkg/modules"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: 5 * time.Second,
	}
}

func (r *Repository) GetUsers(ctx context.Context) ([]modules.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var users []modules.User
	err := r.db.DB.SelectContext(ctx, &users,
		`SELECT id, name, email, age, created_at FROM users ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*modules.User, error) {
	if id <= 0 {
		return nil, apperr.ErrInvalidInput
	}

	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var u modules.User
	err := r.db.DB.GetContext(ctx, &u,
		`SELECT id, name, email, age, created_at FROM users WHERE id = $1`,
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperr.ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) CreateUser(ctx context.Context, in modules.CreateUserInput) (int, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.TrimSpace(in.Email)

	if name == "" || email == "" || in.Age < 0 {
		return 0, apperr.ErrInvalidInput
	}

	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var id int
	err := r.db.DB.QueryRowxContext(
		ctx,
		`INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id`,
		name, email, in.Age,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id int, in modules.UpdateUserInput) error {
	if id <= 0 {
		return apperr.ErrInvalidInput
	}
	if in.Name == nil && in.Email == nil && in.Age == nil {
		return apperr.ErrInvalidInput
	}
	if in.Age != nil && *in.Age < 0 {
		return apperr.ErrInvalidInput
	}

	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	res, err := r.db.DB.ExecContext(ctx, `
		UPDATE users
		SET
			name  = COALESCE($2, name),
			email = COALESCE($3, email),
			age   = COALESCE($4, age)
		WHERE id = $1
	`, id, in.Name, in.Email, in.Age)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.ErrNotFound
	}

	return nil
}

func (r *Repository) DeleteUserByID(ctx context.Context, id int) (int64, error) {
	if id <= 0 {
		return 0, apperr.ErrInvalidInput
	}

	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	res, err := r.db.DB.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if affected == 0 {
		return 0, apperr.ErrNotFound
	}

	return affected, nil
}