
package repository

import (
	"context"
	"os"
	"testing"
)

// TestUserRepository checks actual PostgreSQL user operations.
func TestUserRepository(t *testing.T) {
	// Read the database URL from the environment.
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not configured")
	}

	// Create and verify a PostgreSQL connection pool.
	pool, err := NewPostgresPool(databaseURL)
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer pool.Close()

	// Create the user repository.
	userRepo := NewUserRepository(pool)

	ctx := context.Background()
	phoneNumber := "9876543210"

	// Remove any previous test record to keep this test repeatable.
	_, err = pool.Exec(ctx,
		"DELETE FROM users WHERE phone_number = $1",
		phoneNumber,
	)
	if err != nil {
		t.Fatalf("failed to clean up test user: %v", err)
	}

	// Insert a new user.
	createdUser, err := userRepo.CreateUser(ctx, phoneNumber)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if createdUser.ID == "" {
		t.Fatal("expected generated user ID")
	}

	// Retrieve the user using the phone number.
	foundUser, err := userRepo.GetUserByPhone(ctx, phoneNumber)
	if err != nil {
		t.Fatalf("failed to retrieve user: %v", err)
	}

	// Verify that the retrieved user matches the created user.
	if foundUser.ID != createdUser.ID {
		t.Errorf("user ID mismatch: got %s, want %s",
			foundUser.ID, createdUser.ID)
	}

	if foundUser.PhoneNumber != phoneNumber {
		t.Errorf("phone number mismatch: got %s, want %s",
			foundUser.PhoneNumber, phoneNumber)
	}

	// Clean up the test record.
	_, err = pool.Exec(ctx,
		"DELETE FROM users WHERE phone_number = $1",
		phoneNumber,
	)
	if err != nil {
		t.Errorf("failed to clean up test user: %v", err)
	}
}
