package client

import "testing"

// TestDialServiceRejectsEmptyAddress checks invalid input handling.
func TestDialServiceRejectsEmptyAddress(t *testing.T) {
	// Attempt to create a connection with an empty address.
	conn, err := DialService("")

	// An empty address should return an error.
	if err == nil {
		t.Fatal("expected an error for empty address")
	}

	// No connection should be returned when input is invalid.
	if conn != nil {
		t.Fatal("expected nil connection for empty address")
	}
}
