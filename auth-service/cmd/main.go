package main

import (
	// Import the Auth handler implementation.
	"github.com/rosha/mini-auth-system/auth-service/internal/handler"

	// Import the PostgreSQL repository.
	"github.com/rosha/mini-auth-system/auth-service/internal/repository"

	// Import the Auth gRPC server.
	authserver "github.com/rosha/mini-auth-system/auth-service/internal/server"

	// Import shared configuration and logger.
	"github.com/rosha/mini-auth-system/service-utils/config"
	"github.com/rosha/mini-auth-system/service-utils/logger"
)

func main() {
	// Load Auth service configuration.
	cfg, err := config.LoadAuth()
	if err != nil {
		panic(err)
	}

	// Initialize structured logging.
	log := logger.New()

	// Establish a connection pool to PostgreSQL.
	dbPool, err := repository.NewPostgresPool(cfg.DatabaseURL)
	if err != nil {
		log.Error("Failed to connect to PostgreSQL", "error", err)
		panic(err)
	}

	// Close the database pool when the application exits.
	defer dbPool.Close()

	log.Info("Connected to PostgreSQL successfully")

	// Create the user repository using the PostgreSQL connection pool.
	userRepo := repository.NewUserRepository(dbPool)

	// Create the Auth handler and inject its dependencies.
	authHandler := handler.NewAuthHandler(cfg.JWTSecret, userRepo)

	// Start the gRPC server with the handler registered.
	log.Info("Starting Auth gRPC server", "port", cfg.GRPCPort)

	if err := authserver.StartGRPCServer(cfg.GRPCPort, authHandler); err != nil {
		log.Error("Auth gRPC server stopped", "error", err)
		panic(err)
	}
}
