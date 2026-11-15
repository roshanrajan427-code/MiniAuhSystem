package service

import (
	"testing"
)

// TestJWTGeneration checks that a JWT is generated successfully.
func TestJWTGeneration(t *testing.T) {
	jwtService := NewJWTService("test-secret-key")

	token, err := jwtService.Generate("9876543210")
	if err != nil {
		t.Fatalf("JWT generation failed: %v", err)
	}

	if token == "" {
		t.Fatal("expected a JWT token, got an empty string")
	}
}

// TestJWTValidation checks that a generated JWT can be validated.
func TestJWTValidation(t *testing.T) {
	jwtService := NewJWTService("test-secret-key")

	token, err := jwtService.Generate("9876543210")
	if err != nil {
		t.Fatalf("JWT generation failed: %v", err)
	}

	phoneNumber, err := jwtService.Validate(token)
	if err != nil {
		t.Fatalf("JWT validation failed: %v", err)
	}

	if phoneNumber != "9876543210" {
		t.Errorf("expected phone number 9876543210, got %s", phoneNumber)
	}
}

// TestJWTInvalidSecret checks that a token signed with another secret is rejected.
func TestJWTInvalidSecret(t *testing.T) {
	jwtService := NewJWTService("test-secret-key")
	anotherService := NewJWTService("another-secret-key")

	token, err := jwtService.Generate("9876543210")
	if err != nil {
		t.Fatalf("JWT generation failed: %v", err)
	}

	if _, err := anotherService.Validate(token); err == nil {
		t.Error("expected token validation to fail with an incorrect secret")
	}
}
