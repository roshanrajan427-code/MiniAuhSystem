package config

import (
	"fmt"
	"os"
)

// Config holds configuration required by our microservices.
type Config struct {
	HTTPPort           string
	GRPCPort           string
	AuthServiceAddress string
	DatabaseURL        string
	JWTSecret          string
}

// LoadGateway loads configuration required by Gateway.
func LoadGateway() *Config {
	return &Config{
		HTTPPort:           getEnv("HTTP_PORT", "8080"),
		AuthServiceAddress: getEnv("AUTH_SERVICE_ADDRESS", "localhost:9090"),
	}
}

// LoadAuth loads configuration required by Auth.
func LoadAuth() (*Config, error) {
	cfg := &Config{
		GRPCPort:    getEnv("GRPC_PORT", "9090"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	// Auth cannot start without its database connection.
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	// Auth needs a secret to sign JWT tokens.
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

// getEnv returns the environment variable or a fallback value.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
