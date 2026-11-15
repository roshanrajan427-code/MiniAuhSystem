package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// OTPService manages temporary OTP generation and verification.
type OTPService struct {
	mu      sync.Mutex
	otps    map[string]string
	expires map[string]time.Time
}

// NewOTPService initializes the OTP service.
func NewOTPService() *OTPService {
	return &OTPService{
		otps:    make(map[string]string),
		expires: make(map[string]time.Time),
	}
}

// Generate creates a secure four-digit OTP.
func (s *OTPService) Generate(phoneNumber string) (string, error) {
	// Generate a cryptographically secure number between 0000 and 9999.
	number, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Keep leading zeros so every OTP contains exactly four digits.
	otp := fmt.Sprintf("%04d", number.Int64())

	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the OTP with a five-minute expiration.
	s.otps[phoneNumber] = otp
	s.expires[phoneNumber] = time.Now().Add(5 * time.Minute)

	return otp, nil
}

// Verify checks whether the supplied OTP is valid.
func (s *OTPService) Verify(phoneNumber, submittedOTP string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Retrieve the stored OTP and its expiration time.
	storedOTP, exists := s.otps[phoneNumber]
	expiry := s.expires[phoneNumber]

	// Reject missing or expired OTPs.
	if !exists || time.Now().After(expiry) {
		delete(s.otps, phoneNumber)
		delete(s.expires, phoneNumber)
		return false
	}

	// Reject an incorrect OTP.
	if storedOTP != submittedOTP {
		return false
	}

	// Remove the OTP after successful verification.
	delete(s.otps, phoneNumber)
	delete(s.expires, phoneNumber)

	return true
}

// GetOTPForTesting returns the stored OTP for development tests only.
func (s *OTPService) GetOTPForTesting(phoneNumber string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Retrieve the OTP without exposing the internal map.
	otp, exists := s.otps[phoneNumber]
	return otp, exists
}
