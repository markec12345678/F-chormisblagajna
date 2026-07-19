package config

import (
	"testing"
)

func TestGenerateRandomSecret(t *testing.T) {
	secret1 := generateRandomSecret()
	secret2 := generateRandomSecret()

	if secret1 == "" {
		t.Error("generateRandomSecret returned empty string")
	}

	if secret2 == "" {
		t.Error("generateRandomSecret returned empty string")
	}

	if secret1 == secret2 {
		t.Error("generateRandomSecret should return different values")
	}

	if len(secret1) != 64 {
		t.Errorf("Secret length = %v, want 64 (hex encoded 32 bytes)", len(secret1))
	}
}
