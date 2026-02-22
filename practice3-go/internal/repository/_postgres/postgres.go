package _postgres

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"practice3-go/pkg/modules"
)

type Dialect struct {
	DB *sqlx.DB
}

func NewPGXDialect(ctx context.Context, cfg *modules.PostgreConfig) (*Dialect, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	if err := AutoMigrate(cfg); err != nil {
		return nil, err
	}

	return &Dialect{DB: db}, nil
}

func AutoMigrate(cfg *modules.PostgreConfig) error {
	sourceURL := "file://database/migrations"
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}