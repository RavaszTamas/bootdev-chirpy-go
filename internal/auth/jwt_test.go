package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "mySuperSecretKey123!"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	if token == "" {
		t.Error("Expected token to not be empty")
	}

	// JWTs should be composed of 3 base64url encoded parts separated by dots
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Expected token to have 3 parts separated by dots, got: %d parts", len(parts))
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "mySuperSecretKey123!"

	// Generate a valid token to test against
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to setup token for testing: %v", err)
	}

	t.Run("Valid Token", func(t *testing.T) {
		extractedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Errorf("Unexpected error for valid input: %v", err)
		}
		if extractedID != userID {
			t.Errorf("Expected extracted ID %v, got %v", userID, extractedID)
		}
	})

	t.Run("Invalid Secret", func(t *testing.T) {
		extractedID, err := ValidateJWT(token, "wrongSecret!99")
		if err == nil {
			t.Error("Expected an error when using an incorrect secret, got nil")
		}
		if extractedID != uuid.Nil {
			t.Error("Expected extracted ID to be nil UUID for an incorrect secret")
		}
	})

	t.Run("Malformed Token", func(t *testing.T) {
		extractedID, err := ValidateJWT("not_a_real_token_string", secret)
		if err == nil {
			t.Error("Expected an error when parsing an invalid token string, got nil")
		}
		if extractedID != uuid.Nil {
			t.Error("Expected extracted ID to be nil UUID when using an invalid token")
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
		// Create a token with a negative duration so it is instantly expired
		expiredToken, err := MakeJWT(userID, secret, -1*time.Second)
		if err != nil {
			t.Fatalf("Failed to setup expired token for testing: %v", err)
		}

		extractedID, err := ValidateJWT(expiredToken, secret)
		if err == nil {
			t.Error("Expected an error when validating an expired token, got nil")
		}
		if extractedID != uuid.Nil {
			t.Error("Expected extracted ID to be nil UUID when using an expired token")
		}
	})
}
