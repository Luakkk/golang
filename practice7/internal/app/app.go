package app

import (
	"fmt"
	"log"
	"os"

	"practice-7/config"
	v1 "practice-7/internal/controller/http/v1"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/logger"
	"practice-7/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func Run() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config.New:", err)
	}

	// Set JWT secret for utils package
	os.Setenv("JWT_SECRET", cfg.JWTSecret)

	// Logger
	l := logger.New()

	// Postgres
	pg, err := postgres.New(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	if err != nil {
		log.Fatal("postgres.New:", err)
	}

	// Enable uuid-ossp extension and auto-migrate
	pg.Conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err := pg.Conn.AutoMigrate(&entity.User{}); err != nil {
		log.Fatal("AutoMigrate:", err)
	}

	// Layers
	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Gin router
	router := gin.Default()
	v1.NewRouter(router, userUseCase, l)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	l.Info("Server started on", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("router.Run:", err)
	}
}
