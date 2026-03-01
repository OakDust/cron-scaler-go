package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort       string
	GRPCServerAddr string
	CORSAllowOrigin string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		HTTPPort:        getEnv("HTTP_PORT", "8080"),
		GRPCServerAddr:  getEnv("GRPC_SERVER_ADDR", "localhost:50051"),
		CORSAllowOrigin: getEnv("CORS_ALLOW_ORIGIN", "*"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
