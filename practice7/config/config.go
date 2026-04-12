package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	ServerPort string
	JWTSecret  string
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "auth_service"),
		DBPort:     getEnv("DB_PORT", "5432"),
		ServerPort: getEnv("SERVER_PORT", "8090"),
		JWTSecret:  getEnv("JWT_SECRET", "secret"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
