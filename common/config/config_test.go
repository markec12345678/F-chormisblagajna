package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nutrixpos/pos/common/logger"
)

func testLogger() logger.ILogger {
	l := logger.NewZeroLog()
	return &l
}

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

func TestConfigFactory_UnknownType(t *testing.T) {
	log := testLogger()
	result := ConfigFactory("unknown", "", log)
	if result.Env != "" {
		t.Errorf("unknown type should return empty Config, got Env=%v", result.Env)
	}
}

func TestConfigFactory_Viper(t *testing.T) {
	yaml := `
env: test
auth:
  enabled: true
  jwt_secret: my-secret
  jwt_expire_hrs: 48
databases:
  - host: localhost
    port: 27017
    database: testdb
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(yaml), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	log := testLogger()
	result := ConfigFactory("viper", configPath, log)

	if result.Env != "test" {
		t.Errorf("Env = %v, want 'test'", result.Env)
	}
	if !result.Auth.Enabled {
		t.Error("Auth.Enabled should be true")
	}
	if result.Auth.JWTSecret != "my-secret" {
		t.Errorf("Auth.JWTSecret = %v, want 'my-secret'", result.Auth.JWTSecret)
	}
	if result.Auth.JWTExpireHrs != 48 {
		t.Errorf("Auth.JWTExpireHrs = %v, want 48", result.Auth.JWTExpireHrs)
	}
	if len(result.Databases) != 1 {
		t.Fatalf("Databases length = %v, want 1", len(result.Databases))
	}
	if result.Databases[0].Host != "localhost" {
		t.Errorf("Database Host = %v, want 'localhost'", result.Databases[0].Host)
	}
	if result.Databases[0].Database != "testdb" {
		t.Errorf("Database Name = %v, want 'testdb'", result.Databases[0].Database)
	}
}

func TestConfigFactory_Viper_AutoGeneratesJWTSecret(t *testing.T) {
	yaml := `
env: prod
auth:
  enabled: true
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(yaml), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	log := testLogger()
	result := ConfigFactory("viper", configPath, log)

	if result.Auth.JWTSecret == "" {
		t.Error("JWTSecret should be auto-generated when auth is enabled but secret is empty")
	}
	if len(result.Auth.JWTSecret) != 64 {
		t.Errorf("Auto-generated JWTSecret length = %v, want 64", len(result.Auth.JWTSecret))
	}
}

func TestConfigFactory_Viper_NonexistentFile(t *testing.T) {
	log := testLogger()
	result := ConfigFactory("viper", "/nonexistent/path.yaml", log)

	if result.Env != "" {
		t.Errorf("nonexistent file should return empty Config, got Env=%v", result.Env)
	}
}
