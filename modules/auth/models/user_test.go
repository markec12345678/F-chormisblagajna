package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestUser_ToResponse(t *testing.T) {
	objID := bson.NewObjectID()
	user := User{
		ID:           objID,
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: "secret-hash",
		Roles:        []string{"admin", "cashier"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	response := user.ToResponse()

	if response.ID != objID.Hex() {
		t.Errorf("ID = %v, want %v", response.ID, objID.Hex())
	}
	if response.Username != "admin" {
		t.Errorf("Username = %v, want 'admin'", response.Username)
	}
	if response.Email != "admin@example.com" {
		t.Errorf("Email = %v, want 'admin@example.com'", response.Email)
	}
	if len(response.Roles) != 2 {
		t.Errorf("Roles length = %v, want 2", len(response.Roles))
	}
	if response.Roles[0] != "admin" {
		t.Errorf("Roles[0] = %v, want 'admin'", response.Roles[0])
	}
	if response.Roles[1] != "cashier" {
		t.Errorf("Roles[1] = %v, want 'cashier'", response.Roles[1])
	}
}

func TestUser_ToResponse_EmptyRoles(t *testing.T) {
	objID := bson.NewObjectID()
	user := User{
		ID:       objID,
		Username: "test",
		Email:    "test@example.com",
		Roles:    []string{},
	}

	response := user.ToResponse()

	if response.ID != objID.Hex() {
		t.Errorf("ID = %v, want %v", response.ID, objID.Hex())
	}
	if len(response.Roles) != 0 {
		t.Errorf("Roles length = %v, want 0", len(response.Roles))
	}
}

func TestUserResponse_JSON(t *testing.T) {
	objID := bson.NewObjectID()
	user := User{
		ID:       objID,
		Username: "chef",
		Email:    "chef@example.com",
		Roles:    []string{"chef"},
	}

	response := user.ToResponse()

	if response.Username != "chef" {
		t.Errorf("Username = %v, want 'chef'", response.Username)
	}
	if response.Email != "chef@example.com" {
		t.Errorf("Email = %v, want 'chef@example.com'", response.Email)
	}
}

func TestUser_BsonTags(t *testing.T) {
	user := User{}

	data, err := bson.Marshal(user)
	if err != nil {
		t.Fatalf("BSON marshal failed: %v", err)
	}

	var decoded User
	if err := bson.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("BSON unmarshal failed: %v", err)
	}

	if decoded.Username != "" {
		t.Errorf("Username should be empty after round-trip, got %q", decoded.Username)
	}
}
