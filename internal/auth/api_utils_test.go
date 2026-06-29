package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {

	t.Run("Valid bearer", func(t *testing.T) {
		header := http.Header{}
		expectedToken := "my-secret-token"
		header.Add("Authorization", "Bearer "+expectedToken)

		value, err := GetBearerToken(header)
		if err != nil {
			t.Errorf("Unexpected error for valid header: %v", err)
		}
		if value != expectedToken {
			t.Errorf("Expected token %q, got %q", expectedToken, value)
		}
	})

	t.Run("Missing header", func(t *testing.T) {
		header := http.Header{} // Empty header

		value, err := GetBearerToken(header)
		if err == nil {
			t.Error("Expected an error for missing Authorization header, got nil")
		}
		if value != "" {
			t.Errorf("Expected empty string for missing header, got %q", value)
		}
	})

	t.Run("Empty header value", func(t *testing.T) {
		header := http.Header{}
		header.Add("Authorization", "")

		value, err := GetBearerToken(header)
		if err == nil {
			t.Error("Expected an error for empty Authorization header value, got nil")
		}
		if value != "" {
			t.Errorf("Expected empty string for empty header value, got %q", value)
		}
	})

	t.Run("Token with extra whitespace", func(t *testing.T) {
		header := http.Header{}
		expectedToken := "my-secret-token"
		// Simulating sloppy formatting from a client
		header.Add("Authorization", "Bearer     "+expectedToken+"   ")

		value, err := GetBearerToken(header)
		if err != nil {
			t.Errorf("Unexpected error for valid header with extra spaces: %v", err)
		}
		if value != expectedToken {
			t.Errorf("Expected token %q, got %q", expectedToken, value)
		}
	})
}
