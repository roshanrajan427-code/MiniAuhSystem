
package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/rosha/mini-auth-system/auth-service/internal/repository"
	"github.com/rosha/mini-auth-system/auth-service/internal/service"
	authv1 "github.com/rosha/mini-auth-system/protobufs/generated/auth/v1"
)

// AuthHandler implements the generated AuthServiceServer interface.
type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer

	// Services handle OTP and JWT business logic.
	otpService *service.OTPService
	jwtService *service.JWTService

	// Repository handles PostgreSQL user operations.
	userRepo *repository.UserRepository
}

// NewAuthHandler creates a handler with its required dependencies.
func NewAuthHandler(
	jwtSecret string,
	userRepo *repository.UserRepository,
) *AuthHandler {
	return &AuthHandler{
		otpService: service.NewOTPService(),
		jwtService: service.NewJWTService(jwtSecret),
		userRepo:   userRepo,
	}
}

// Login generates an OTP for the supplied phone number.
func (h *AuthHandler) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {

	// Reject an empty phone number.
	if req.GetPhoneNumber() == "" {
		return nil, fmt.Errorf("phone number is required")
	}

	// Generate and store the OTP.
	otp, err := h.otpService.Generate(req.GetPhoneNumber())
	if err != nil {
		return nil, fmt.Errorf("OTP generation failed: %w", err)
	}

	// Print OTP only for local development.
	fmt.Printf("Development OTP for %s: %s\n", req.GetPhoneNumber(), otp)

	return &authv1.LoginResponse{
		Message: "OTP generated successfully",
	}, nil
}

// VerifyOTP verifies the OTP, creates or retrieves the user, and generates a JWT.
func (h *AuthHandler) VerifyOTP(
	ctx context.Context,
	req *authv1.VerifyOTPRequest,
) (*authv1.VerifyOTPResponse, error) {

	// Reject missing request values.
	if req.GetPhoneNumber() == "" || req.GetOtp() == "" {
		return nil, fmt.Errorf("phone number and OTP are required")
	}

	// Verify the OTP associated with the phone number.
	if !h.otpService.Verify(req.GetPhoneNumber(), req.GetOtp()) {
		return nil, fmt.Errorf("invalid or expired OTP")
	}

	// Check whether this phone number already belongs to a user.
	user, err := h.userRepo.GetUserByPhone(ctx, req.GetPhoneNumber())

	if errors.Is(err, pgx.ErrNoRows) {
		// Create a user only when no matching record exists.
		user, err = h.userRepo.CreateUser(ctx, req.GetPhoneNumber())
		if err != nil {
			return nil, fmt.Errorf("user creation failed: %w", err)
		}
	} else if err != nil {
		// Return actual database errors instead of treating them as missing users.
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	// Generate an access token after successful verification.
	token, err := h.jwtService.Generate(user.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("JWT generation failed: %w", err)
	}

	return &authv1.VerifyOTPResponse{
		AccessToken: token,
	}, nil
}

// GetProfile validates the JWT and returns the user's database profile.
func (h *AuthHandler) GetProfile(
	ctx context.Context,
	req *authv1.GetProfileRequest,
) (*authv1.GetProfileResponse, error) {

	// Reject an empty access token.
	if req.GetAccessToken() == "" {
		return nil, fmt.Errorf("access token is required")
	}

	// Validate the JWT and extract the phone number.
	phoneNumber, err := h.jwtService.Validate(req.GetAccessToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Retrieve the actual user record from PostgreSQL.
	user, err := h.userRepo.GetUserByPhone(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user profile: %w", err)
	}

	return &authv1.GetProfileResponse{
		UserId:      user.ID,
		PhoneNumber: user.PhoneNumber,
	}, nil
}
