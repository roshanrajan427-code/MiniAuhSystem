package server

import (
	"fmt"

	"google.golang.org/grpc"
)

// NewGRPCServer creates a reusable gRPC server instance.
func NewGRPCServer() (*grpc.Server, error) {
	// Create a new gRPC server with default settings.
	grpcServer := grpc.NewServer()

	// Return the server so the Auth service can register its RPC handlers.
	if grpcServer == nil {
		return nil, fmt.Errorf("failed to create gRPC server")
	}

	return grpcServer, nil
}
