package server

import (
	"fmt"
	"net"

	// Import the generated Auth gRPC service.
	authv1 "github.com/rosha/mini-auth-system/protobufs/generated/auth/v1"

	// Import the shared gRPC server utility.
	grpcutil "github.com/rosha/mini-auth-system/service-utils/server"
)

// StartGRPCServer starts the Auth service gRPC server.
func StartGRPCServer(port string, authHandler authv1.AuthServiceServer) error {

	// Open a TCP listener on the configured port.
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	// Create the gRPC server using our shared utility.
	grpcServer, err := grpcutil.NewGRPCServer()
	if err != nil {
		return fmt.Errorf("failed to create gRPC server: %w", err)
	}

	// Register the Auth service implementation.
	authv1.RegisterAuthServiceServer(grpcServer, authHandler)

	// Start accepting incoming gRPC requests.
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("gRPC server failed: %w", err)
	}

	return nil
}
