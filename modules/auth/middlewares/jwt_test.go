package middlewares

import (
	"testing"
	"time"

	"github.com/nutrixpos/pos/modules/auth/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestJWTUtil_GenerateToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key", 24)

	user := models.User{
		ID:       bson.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"admin"},
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("GenerateToken returned empty token")
	}
}

func TestJWTUtil_ValidateToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key", 24)

	user := models.User{
		ID:       bson.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"admin"},
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := jwtUtil.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.Username != user.Username {
		t.Errorf("Username = %v, want %v", claims.Username, user.Username)
	}

	if claims.Email != user.Email {
		t.Errorf("Email = %v, want %v", claims.Email, user.Email)
	}

	if len(claims.Roles) != len(user.Roles) {
		t.Errorf("Roles length = %v, want %v", len(claims.Roles), len(user.Roles))
	}
}

func TestJWTUtil_ValidateToken_InvalidSecret(t *testing.T) {
	jwtUtil1 := NewJWTUtil("secret-1", 24)
	jwtUtil2 := NewJWTUtil("secret-2", 24)

	user := models.User{
		ID:       bson.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"admin"},
	}

	token, err := jwtUtil1.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = jwtUtil2.ValidateToken(token)
	if err != ErrInvalidToken {
		t.Errorf("ValidateToken with wrong secret should return ErrInvalidToken, got: %v", err)
	}
}

func TestJWTUtil_ValidateToken_Expired(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key", -1)

	user := models.User{
		ID:       bson.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"admin"},
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	_, err = jwtUtil.ValidateToken(token)
	if err != ErrExpiredToken {
		t.Errorf("ValidateToken with expired token should return ErrExpiredToken, got: %v", err)
	}
}

func TestJWTUtil_RefreshToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key", 24)

	user := models.User{
		ID:       bson.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"admin"},
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	newToken, err := jwtUtil.RefreshToken(token)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}

	if newToken == "" {
		t.Fatal("RefreshToken returned empty token")
	}

	claims, err := jwtUtil.ValidateToken(newToken)
	if err != nil {
		t.Fatalf("ValidateToken on refreshed token failed: %v", err)
	}

	if claims.Username != user.Username {
		t.Errorf("Username = %v, want %v", claims.Username, user.Username)
	}
}

func TestJWTUtil_RefreshToken_InvalidToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key", 24)

	_, err := jwtUtil.RefreshToken("invalid.token.here")
	if err != ErrInvalidToken {
		t.Errorf("RefreshToken with invalid token should return ErrInvalidToken, got: %v", err)
	}
}

func TestJWTUtil_Roles(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key", 24)

	user := models.User{
		ID:       bson.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"admin", "cashier"},
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := jwtUtil.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if len(claims.Roles) != 2 {
		t.Errorf("Roles length = %v, want 2", len(claims.Roles))
	}
	if claims.Roles[0] != "admin" {
		t.Errorf("Roles[0] = %v, want admin", claims.Roles[0])
	}
	if claims.Roles[1] != "cashier" {
		t.Errorf("Roles[1] = %v, want cashier", claims.Roles[1])
	}
}

func TestNewJWTUtil(t *testing.T) {
	jwtUtil := NewJWTUtil("my-secret", 48)
	if string(jwtUtil.Secret) != "my-secret" {
		t.Errorf("Secret = %v, want 'my-secret'", string(jwtUtil.Secret))
	}
	if jwtUtil.ExpireHours != 48 {
		t.Errorf("ExpireHours = %v, want 48", jwtUtil.ExpireHours)
	}
}

func TestMustObjectIDFromHex_Valid(t *testing.T) {
	oid := bson.NewObjectID()
	result := mustObjectIDFromHex(oid.Hex())
	if result != oid {
		t.Errorf("mustObjectIDFromHex(valid) = %v, want %v", result, oid)
	}
}

func TestMustObjectIDFromHex_Invalid(t *testing.T) {
	result := mustObjectIDFromHex("not-a-valid-hex")
	if result != bson.NilObjectID {
		t.Errorf("mustObjectIDFromHex(invalid) = %v, want NilObjectID", result)
	}
}

func TestMustObjectIDFromHex_Empty(t *testing.T) {
	result := mustObjectIDFromHex("")
	if result != bson.NilObjectID {
		t.Errorf("mustObjectIDFromHex(empty) = %v, want NilObjectID", result)
	}
}

func TestJWTUtil_ValidateToken_NonHMAC(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key", 24)

	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzIn0.fakesignature"

	_, err := jwtUtil.ValidateToken(tokenString)
	if err != ErrInvalidToken {
		t.Errorf("ValidateToken with non-HMAC token should return ErrInvalidToken, got: %v", err)
	}
}

func TestJWTUtil_ValidateToken_Malformed(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key", 24)

	_, err := jwtUtil.ValidateToken("completely.malformed")
	if err != ErrInvalidToken {
		t.Errorf("ValidateToken with malformed token should return ErrInvalidToken, got: %v", err)
	}
}

func TestJWTUtil_GenerateToken_ClaimsPreserved(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key", 24)
	objID := bson.NewObjectID()

	user := models.User{
		ID:       objID,
		Username: "chef1",
		Email:    "chef@kitchen.com",
		Roles:    []string{"chef", "cashier"},
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := jwtUtil.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != objID.Hex() {
		t.Errorf("UserID = %v, want %v", claims.UserID, objID.Hex())
	}
	if claims.Username != "chef1" {
		t.Errorf("Username = %v, want 'chef1'", claims.Username)
	}
	if claims.Email != "chef@kitchen.com" {
		t.Errorf("Email = %v, want 'chef@kitchen.com'", claims.Email)
	}
	if len(claims.Roles) != 2 {
		t.Errorf("Roles length = %v, want 2", len(claims.Roles))
	}
	if claims.Issuer != "nutrix-pos" {
		t.Errorf("Issuer = %v, want 'nutrix-pos'", claims.Issuer)
	}
}

func TestJWTUtil_RefreshToken_RoundTrip(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key", 24)
	objID := bson.NewObjectID()

	user := models.User{
		ID:       objID,
		Username: "roundtrip",
		Email:    "rt@test.com",
		Roles:    []string{"admin"},
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	refreshed, err := jwtUtil.RefreshToken(token)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}

	claims, err := jwtUtil.ValidateToken(refreshed)
	if err != nil {
		t.Fatalf("ValidateToken on refreshed token failed: %v", err)
	}

	if claims.UserID != objID.Hex() {
		t.Errorf("Refreshed UserID = %v, want %v", claims.UserID, objID.Hex())
	}
	if claims.Username != "roundtrip" {
		t.Errorf("Refreshed Username = %v, want 'roundtrip'", claims.Username)
	}
}
