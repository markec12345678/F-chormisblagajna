package middlewares

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("testpassword")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}

	if hash == "testpassword" {
		t.Fatal("HashPassword should not return plaintext")
	}
}

func TestCheckPassword_Correct(t *testing.T) {
	password := "mypassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for correct password")
	}
}

func TestCheckPassword_Incorrect(t *testing.T) {
	hash, err := HashPassword("correctpassword")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if CheckPassword("wrongpassword", hash) {
		t.Error("CheckPassword should return false for wrong password")
	}
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if CheckPassword("", hash) {
		t.Error("CheckPassword should return false for empty password")
	}
}

func TestCheckPassword_EmptyHash(t *testing.T) {
	if CheckPassword("password", "") {
		t.Error("CheckPassword should return false for empty hash")
	}
}

func TestHashPassword_SaltUniqueness(t *testing.T) {
	hash1, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	hash2, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash1 == hash2 {
		t.Error("Two hashes of the same password should differ (different salts)")
	}
	if !CheckPassword("same-password", hash1) {
		t.Error("CheckPassword should return true for first hash")
	}
	if !CheckPassword("same-password", hash2) {
		t.Error("CheckPassword should return true for second hash")
	}
}

func TestCheckPassword_Unicode(t *testing.T) {
	password := "šifra-ñ-日本語-🔑"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for unicode password")
	}
	if CheckPassword("šifra-ñ-日本語-❌", hash) {
		t.Error("CheckPassword should return false for different unicode password")
	}
}

func TestCheckPassword_LongPassword(t *testing.T) {
	longPass := "a]veryLongButUnder72BytesPasswordThatIsStillQuiteLongAndShouldWork"
	hash, err := HashPassword(longPass)
	if err != nil {
		t.Fatalf("HashPassword failed for long password: %v", err)
	}
	if !CheckPassword(longPass, hash) {
		t.Error("CheckPassword should return true for long password")
	}
}

func TestCheckPassword_MaxLengthPassword(t *testing.T) {
	maxPass := ""
	for i := 0; i < 72; i++ {
		maxPass += "x"
	}
	hash, err := HashPassword(maxPass)
	if err != nil {
		t.Fatalf("HashPassword failed for 72-char password: %v", err)
	}
	if !CheckPassword(maxPass, hash) {
		t.Error("CheckPassword should return true for 72-char password")
	}
}

func TestCheckPassword_SpecialChars(t *testing.T) {
	password := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for special chars password")
	}
}
