package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService manages JWT creation and validation.
type JWTService struct {
	secret []byte
}

// NewJWTService creates a JWT service using the configured secret.
func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret: []byte(secret),
	}
}

// Generate creates a signed JWT for the authenticated user.
func (s *JWTService) Generate(phoneNumber string) (string, error) {
	// Define the information stored inside the JWT.
	claims := jwt.MapClaims{
		"phone_number": phoneNumber,
		"iat":          time.Now().Unix(),
		"exp":          time.Now().Add(24 * time.Hour).Unix(),
	}

	// Create a token using the HS256 signing algorithm.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using our secret key.
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signedToken, nil
}

// Validate checks the signature and expiration of a JWT.
func (s *JWTService) Validate(tokenString string) (string, error) {
	// Parse and verify the token.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Reject tokens that don't use the expected signing algorithm.
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return s.secret, nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid JWT: %w", err)
	}

	// Ensure the token is valid and contains the expected claims.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid JWT claims")
	}

	// Extract the phone number from the token.
	phoneNumber, ok := claims["phone_number"].(string)
	if !ok || phoneNumber == "" {
		return "", fmt.Errorf("phone number missing from JWT")
	}

	return phoneNumber, nil
}
