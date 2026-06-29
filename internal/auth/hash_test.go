package auth

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mySuperSecretPassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if hash == "" {
		t.Error("Expected hash to not be empty")
	}

	// argon2id hashes should be in the standard encoded format
	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Errorf("Expected hash to start with '$argon2id$', got: %s", hash)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySuperSecretPassword123!"

	// Generate a valid hash to test against
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to setup hash for testing: %v", err)
	}

	t.Run("Valid Password", func(t *testing.T) {
		match, err := CheckPasswordHash(password, hash)
		if err != nil {
			t.Errorf("Unexpected error for valid input: %v", err)
		}
		if !match {
			t.Error("Expected CheckPasswordHash to return true for the correct password")
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		match, err := CheckPasswordHash("wrongPassword!99", hash)
		if err != nil {
			t.Errorf("Unexpected error for incorrect password: %v", err)
		}
		if match {
			t.Error("Expected CheckPasswordHash to return false for an incorrect password")
		}
	})

	t.Run("Malformed Hash", func(t *testing.T) {
		match, err := CheckPasswordHash(password, "not_a_real_hash_string")
		if err == nil {
			t.Error("Expected an error when parsing an invalid hash string, got nil")
		}
		if match {
			t.Error("Expected match to be false when using an invalid hash")
		}
	})
}
