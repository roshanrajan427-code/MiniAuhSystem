
package handler

import (
	"context"
	"os"
	"testing"

	"github.com/rosha/mini-auth-system/auth-service/internal/repository"
	authv1 "github.com/rosha/mini-auth-system/protobufs/generated/auth/v1"
)

// newTestHandler creates a handler connected to PostgreSQL for integration tests.
func newTestHandler(t *testing.T) *AuthHandler {
	t.Helper()

	// Read the database connection URL from the environment.
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not configured")
	}

	// Create a PostgreSQL connection pool.
	pool, err := repository.NewPostgresPool(databaseURL)
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	// Close the pool when the test finishes.
	t.Cleanup(pool.Close)

	// Create the user repository.
	userRepo := repository.NewUserRepository(pool)

	// Inject the repository into the handler.
	return NewAuthHandler("test-secret", userRepo)
}

// TestLogin verifies that a login request generates an OTP.
func TestLogin(t *testing.T) {
	h := newTestHandler(t)

	response, err := h.Login(context.Background(), &authv1.LoginRequest{
		PhoneNumber: "9876543210",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if response.GetMessage() != "OTP generated successfully" {
		t.Errorf("unexpected response: %s", response.GetMessage())
	}
}

// TestVerifyOTP verifies that a correct OTP produces a JWT.
func TestVerifyOTP(t *testing.T) {
	h := newTestHandler(t)

	phoneNumber := "9876543211"

	// Generate an OTP through the OTP service.
	otp, err := h.otpService.Generate(phoneNumber)
	if err != nil {
		t.Fatalf("OTP generation failed: %v", err)
	}

	response, err := h.VerifyOTP(context.Background(), &authv1.VerifyOTPRequest{
		PhoneNumber: phoneNumber,
		Otp:         otp,
	})
	if err != nil {
		t.Fatalf("OTP verification failed: %v", err)
	}

	if response.GetAccessToken() == "" {
		t.Error("expected an access token")
	}
}

// TestVerifyIncorrectOTP verifies that an incorrect OTP is rejected.
func TestVerifyIncorrectOTP(t *testing.T) {
	h := newTestHandler(t)

	phoneNumber := "9876543212"

	// Generate an OTP so the phone number exists in OTP storage.
	_, err := h.otpService.Generate(phoneNumber)
	if err != nil {
		t.Fatalf("OTP generation failed: %v", err)
	}

	_, err = h.VerifyOTP(context.Background(), &authv1.VerifyOTPRequest{
		PhoneNumber: phoneNumber,
		Otp:         "wrong-otp",
	})

	if err == nil {
		t.Error("expected incorrect OTP to be rejected")
	}
}

// TestVerifyOTPReuse verifies that a successfully used OTP cannot be reused.
func TestVerifyOTPReuse(t *testing.T) {
	h := newTestHandler(t)

	phoneNumber := "9876543213"

	otp, err := h.otpService.Generate(phoneNumber)
	if err != nil {
		t.Fatalf("OTP generation failed: %v", err)
	}

	request := &authv1.VerifyOTPRequest{
		PhoneNumber: phoneNumber,
		Otp:         otp,
	}

	// First verification should succeed.
	if _, err := h.VerifyOTP(context.Background(), request); err != nil {
		t.Fatalf("first verification failed: %v", err)
	}

	// Second verification should fail because the OTP was consumed.
	if _, err := h.VerifyOTP(context.Background(), request); err == nil {
		t.Error("expected reused OTP to be rejected")
	}
}

// TestCompleteAuthenticationFlow tests Login, VerifyOTP, and GetProfile.
func TestCompleteAuthenticationFlow(t *testing.T) {
	h := newTestHandler(t)

	ctx := context.Background()
	phoneNumber := "9876543214"

	// Step 1: Generate an OTP.
	loginResponse, err := h.Login(ctx, &authv1.LoginRequest{
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if loginResponse.GetMessage() != "OTP generated successfully" {
		t.Fatalf("unexpected login response: %s", loginResponse.GetMessage())
	}

	// Step 2: Retrieve the OTP using the test-only helper.
	otp, exists := h.otpService.GetOTPForTesting(phoneNumber)
	if !exists {
		t.Fatal("expected OTP to be stored")
	}

	// Step 3: Verify OTP and receive a JWT.
	verifyResponse, err := h.VerifyOTP(ctx, &authv1.VerifyOTPRequest{
		PhoneNumber: phoneNumber,
		Otp:         otp,
	})
	if err != nil {
		t.Fatalf("VerifyOTP failed: %v", err)
	}

	accessToken := verifyResponse.GetAccessToken()
	if accessToken == "" {
		t.Fatal("expected an access token")
	}

	// Step 4: Retrieve the profile using the JWT.
	profileResponse, err := h.GetProfile(ctx, &authv1.GetProfileRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}

	// Step 5: Verify the returned database profile.
	if profileResponse.GetPhoneNumber() != phoneNumber {
		t.Errorf("expected phone number %s, got %s",
			phoneNumber, profileResponse.GetPhoneNumber())
	}

	if profileResponse.GetUserId() == "" {
		t.Error("expected a database-generated user ID")
	}
}
