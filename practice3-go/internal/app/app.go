package app

import (
	"context"
	"os"
	"time"

	"practice3-go/internal/repository"
	"practice3-go/internal/repository/_postgres"
	"practice3-go/internal/usecase"
	"practice3-go/pkg/modules"
)

func LoadConfig() modules.AppConfig {
	pg := modules.PostgreConfig{
		Host:        getenv("PG_HOST", "localhost"),
		Port:        getenv("PG_PORT", "5433"),
		Username:    getenv("PG_USER", "alua"),
		Password:    getenv("PG_PASSWORD", "postgres"),
		DBName:      getenv("PG_DBNAME", "mydb"),
		SSLMode:     getenv("PG_SSLMODE", "disable"),
		ExecTimeout: 5 * time.Second,
	}

	return modules.AppConfig{
		Postgres: pg,
		APIKey:   getenv("API_KEY", "supersecret"),
	}
}

func Build(ctx context.Context, cfg modules.AppConfig) (*_postgres.Dialect, *usecase.UserUsecase, error) {
	db, err := _postgres.NewPGXDialect(ctx, &cfg.Postgres)
	if err != nil {
		return nil, nil, err
	}

	repos := repository.NewRepositories(db)
	uc := usecase.NewUserUsecase(repos.UserRepository)
	return db, uc, nil
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}