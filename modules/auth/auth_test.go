package auth

import (
	"testing"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/modules/core/models"
)

func TestNewBuilder(t *testing.T) {
	conf := config.Config{Env: "test"}
	settings := models.Settings{
		ShopMode: "kitchen",
	}

	mb := NewBuilder(conf, settings)

	if mb == nil {
		t.Fatal("NewBuilder returned nil")
	}
	if mb.Config.Env != "test" {
		t.Errorf("Config.Env = %v, want 'test'", mb.Config.Env)
	}
	if mb.Settings.ShopMode != "kitchen" {
		t.Errorf("Settings.ShopMode = %v, want 'kitchen'", mb.Settings.ShopMode)
	}
}

func TestNewBuilder_EmptyConfig(t *testing.T) {
	mb := NewBuilder(config.Config{}, models.Settings{})

	if mb == nil {
		t.Fatal("NewBuilder with empty config returned nil")
	}
	if mb.Config.Env != "" {
		t.Errorf("Config.Env should be empty, got %v", mb.Config.Env)
	}
}

func TestAuthModuleBuilder_Fields(t *testing.T) {
	conf := config.Config{Env: "prod"}
	settings := models.Settings{ShopMode: "retail"}

	mb := NewBuilder(conf, settings)

	if mb.Logger != nil {
		t.Error("Logger should be nil initially")
	}
	if mb.Prompter != nil {
		t.Error("Prompter should be nil initially")
	}
}
