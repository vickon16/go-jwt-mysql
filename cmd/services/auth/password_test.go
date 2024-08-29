package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Errorf("expected hash to be non-empty")
	}

	if hash == "password" {
		t.Errorf("expected hash to be different from password")
	}

}

func TestComparePasswords(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}

	if !ComparePassword(hash, "password") {
		t.Errorf("expected password to match")
	}

	if ComparePassword(hash, "wrong-password") {
		t.Errorf("expected password to not match")
	}

}
