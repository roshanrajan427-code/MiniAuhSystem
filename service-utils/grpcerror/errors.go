package grpcerror

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// New creates a standardized gRPC status error.
func New(code codes.Code, message string) error {
	// Convert the supplied gRPC code and message into a status error.
	return status.Error(code, message)
}

// InvalidArgument represents invalid client input.
func InvalidArgument(message string) error {
	return New(codes.InvalidArgument, message)
}

// Unauthenticated represents authentication failure.
func Unauthenticated(message string) error {
	return New(codes.Unauthenticated, message)
}

// NotFound represents a missing resource.
func NotFound(message string) error {
	return New(codes.NotFound, message)
}

// Internal represents an unexpected internal error.
func Internal(message string) error {
	return New(codes.Internal, message)
}
