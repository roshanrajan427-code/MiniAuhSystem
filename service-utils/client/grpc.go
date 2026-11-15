package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DialService establishes a connection to a gRPC service.
//
// address is the host and port of the remote gRPC server,
// for example: localhost:9090.
func DialService(address string) (*grpc.ClientConn, error) {
	// Reject an empty address before attempting to connect.
	if address == "" {
		return nil, fmt.Errorf("gRPC service address cannot be empty")
	}

	// Create a gRPC client connection.
	// NewClient creates the connection; it does not guarantee
	// that the remote server is already reachable.
	conn, err := grpc.NewClient(
		address,

		// Use insecure credentials for local development only.
		// We'll configure TLS before production deployment.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	// Return the connection so the calling service can use it.
	return conn, nil
}
