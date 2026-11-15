package service

import (
	"testing"
)

// TestOTPGeneration checks that a generated OTP has exactly four digits.
func TestOTPGeneration(t *testing.T) {
	otpService := NewOTPService()

	otp, err := otpService.Generate("9876543210")
	if err != nil {
		t.Fatalf("OTP generation failed: %v", err)
	}

	if len(otp) != 4 {
		t.Errorf("expected 4-digit OTP, got %s", otp)
	}
}

// TestOTPVerification checks successful OTP verification.
func TestOTPVerification(t *testing.T) {
	otpService := NewOTPService()

	otp, err := otpService.Generate("9876543210")
	if err != nil {
		t.Fatalf("OTP generation failed: %v", err)
	}

	if !otpService.Verify("9876543210", otp) {
		t.Error("expected OTP verification to succeed")
	}

	// An OTP should only be usable once.
	if otpService.Verify("9876543210", otp) {
		t.Error("expected OTP reuse to fail")
	}
}

// TestIncorrectOTP checks that an incorrect OTP is rejected.
func TestIncorrectOTP(t *testing.T) {
	otpService := NewOTPService()

	_, err := otpService.Generate("9876543210")
	if err != nil {
		t.Fatalf("OTP generation failed: %v", err)
	}

	if otpService.Verify("9876543210", "0000") {
		t.Error("expected incorrect OTP verification to fail")
	}
}
