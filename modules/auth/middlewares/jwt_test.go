package middlewares

import (
	"testing"
	"time"

	"github.com/nutrixpos/pos/modules/auth/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJWTUtil_GenerateToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key", 24)

	user := models.User{
		ID:       primitive.NewObjectID(),
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
		ID:       primitive.NewObjectID(),
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
		ID:       primitive.NewObjectID(),
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
		ID:       primitive.NewObjectID(),
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
		ID:       primitive.NewObjectID(),
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
