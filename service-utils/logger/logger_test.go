package logger

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestLoggerProducesJSON(t *testing.T) {
	// Buffer captures log output instead of printing it to the terminal.
	var output bytes.Buffer

	// Create a logger that writes into our buffer.
	log := NewWithOutput(&output)

	// Write a sample informational log.
	log.Info("Login request received",
		"service", "auth-service",
		"phone_number", "+919876543210",
	)

	// Decode the captured JSON output.
	var entry map[string]any

	if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
		t.Fatalf("expected valid JSON log, got error: %v", err)
	}

	// Verify the log message.
	if entry["msg"] != "Login request received" {
		t.Errorf("unexpected log message: %v", entry["msg"])
	}

	// Verify the structured service field.
	if entry["service"] != "auth-service" {
		t.Errorf("unexpected service field: %v", entry["service"])
	}
}
