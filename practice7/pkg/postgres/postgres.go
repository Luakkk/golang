package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	Conn *gorm.DB
}

func New(host, user, password, dbname, port string) (*Postgres, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres.New: %w", err)
	}
	return &Postgres{Conn: db}, nil
}
