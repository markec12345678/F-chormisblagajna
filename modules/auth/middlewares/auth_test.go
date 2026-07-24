package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/modules/auth/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}

func TestNoAuth_AllowAuthenticated(t *testing.T) {
	na := NewNoAuth(config.Config{})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	na.AllowAuthenticated(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("NoAuth.AllowAuthenticated status = %d, want %d", rr.Code, http.StatusOK)
	}
	if rr.Body.String() != "ok" {
		t.Errorf("NoAuth.AllowAuthenticated body = %q, want %q", rr.Body.String(), "ok")
	}
}

func TestNoAuth_AllowAnyOfRoles(t *testing.T) {
	na := NewNoAuth(config.Config{})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	na.AllowAnyOfRoles(okHandler(), "admin").ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("NoAuth.AllowAnyOfRoles status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestInternalAuth_AllowAuthenticated_Disabled(t *testing.T) {
	conf := config.Config{Auth: config.AuthConfig{Enabled: false}}
	ia := NewInternalAuth(conf, nil)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ia.AllowAuthenticated(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("auth disabled should pass through, got status %d", rr.Code)
	}
}

func TestInternalAuth_AllowAuthenticated_NoHeader(t *testing.T) {
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}}
	ia := NewInternalAuth(conf, NewJWTUtil("secret", 24))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ia.AllowAuthenticated(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("no auth header should return 403, got %d", rr.Code)
	}
}

func TestInternalAuth_AllowAuthenticated_InvalidToken(t *testing.T) {
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}, Zitadel: config.ZitadelConfig{Enabled: false}}
	ia := NewInternalAuth(conf, NewJWTUtil("secret", 24))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.value")

	ia.AllowAuthenticated(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("invalid token without zitadel should return 403, got %d", rr.Code)
	}
}

func TestInternalAuth_AllowAuthenticated_InvalidToken_Zitadel(t *testing.T) {
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}, Zitadel: config.ZitadelConfig{Enabled: true}}
	ia := NewInternalAuth(conf, NewJWTUtil("secret", 24))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.value")

	ia.AllowAuthenticated(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("invalid token with zitadel enabled should return 401, got %d", rr.Code)
	}
}

func TestInternalAuth_AllowAuthenticated_ValidToken(t *testing.T) {
	jwtUtil := NewJWTUtil("secret", 24)
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}}
	ia := NewInternalAuth(conf, jwtUtil)

	token, err := jwtUtil.GenerateToken(testUser())
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	ia.AllowAuthenticated(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("valid token should return 200, got %d", rr.Code)
	}
}

func TestInternalAuth_AllowAuthenticated_BearerPrefixStripped(t *testing.T) {
	jwtUtil := NewJWTUtil("secret", 24)
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}}
	ia := NewInternalAuth(conf, jwtUtil)

	token, err := jwtUtil.GenerateToken(testUser())
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)

	ia.AllowAuthenticated(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("valid token without Bearer prefix should still work, got %d", rr.Code)
	}
}

func TestInternalAuth_AllowAuthenticated_EmptyToken(t *testing.T) {
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}}
	ia := NewInternalAuth(conf, NewJWTUtil("secret", 24))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer ")

	ia.AllowAuthenticated(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("empty token should return 403, got %d", rr.Code)
	}
}

func TestInternalAuth_AllowAnyOfRoles_Disabled(t *testing.T) {
	conf := config.Config{Auth: config.AuthConfig{Enabled: false}}
	ia := NewInternalAuth(conf, nil)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ia.AllowAnyOfRoles(okHandler(), "admin").ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("auth disabled should pass through, got status %d", rr.Code)
	}
}

func TestInternalAuth_AllowAnyOfRoles_NoHeader(t *testing.T) {
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}}
	ia := NewInternalAuth(conf, NewJWTUtil("secret", 24))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ia.AllowAnyOfRoles(okHandler(), "admin").ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("no auth header should return 403, got %d", rr.Code)
	}
}

func TestInternalAuth_AllowAnyOfRoles_InvalidToken(t *testing.T) {
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}}
	ia := NewInternalAuth(conf, NewJWTUtil("secret", 24))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")

	ia.AllowAnyOfRoles(okHandler(), "admin").ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("invalid token should return 401, got %d", rr.Code)
	}
}

func TestInternalAuth_AllowAnyOfRoles_SuperuserBypass(t *testing.T) {
	jwtUtil := NewJWTUtil("secret", 24)
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}}
	ia := NewInternalAuth(conf, jwtUtil)

	token, err := jwtUtil.GenerateToken(testUserWithRoles([]string{"superuser"}))
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	ia.AllowAnyOfRoles(okHandler(), "admin").ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("superuser should bypass role check, got %d", rr.Code)
	}
}

func TestInternalAuth_AllowAnyOfRoles_MatchingRole(t *testing.T) {
	jwtUtil := NewJWTUtil("secret", 24)
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}}
	ia := NewInternalAuth(conf, jwtUtil)

	token, err := jwtUtil.GenerateToken(testUserWithRoles([]string{"cashier"}))
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	ia.AllowAnyOfRoles(okHandler(), "cashier", "manager").ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("matching role should return 200, got %d", rr.Code)
	}
}

func TestInternalAuth_AllowAnyOfRoles_NoMatchingRole(t *testing.T) {
	jwtUtil := NewJWTUtil("secret", 24)
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}}
	ia := NewInternalAuth(conf, jwtUtil)

	token, err := jwtUtil.GenerateToken(testUserWithRoles([]string{"viewer"}))
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	ia.AllowAnyOfRoles(okHandler(), "admin", "manager").ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("non-matching role should return 403, got %d", rr.Code)
	}
}

func TestNewInternalAuth(t *testing.T) {
	conf := config.Config{Auth: config.AuthConfig{Enabled: true}}
	jwtUtil := NewJWTUtil("secret", 24)
	ia := NewInternalAuth(conf, jwtUtil)

	if ia.JWTUtil != jwtUtil {
		t.Error("NewInternalAuth should set JWTUtil")
	}
	if ia.Config.Auth.Enabled != true {
		t.Error("NewInternalAuth should set Config")
	}
}

func TestNewNoAuth(t *testing.T) {
	conf := config.Config{Env: "test"}
	na := NewNoAuth(conf)

	if na.Config.Env != "test" {
		t.Errorf("NewNoAuth should set Config, got %v", na.Config.Env)
	}
}

func testUser() models.User {
	return models.User{
		ID:       bson.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"admin"},
	}
}

func testUserWithRoles(roles []string) models.User {
	return models.User{
		ID:       bson.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    roles,
	}
}
