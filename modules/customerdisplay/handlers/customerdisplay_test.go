package handlers

import (
	"encoding/json"
	"testing"

	cd_models "github.com/nutrixpos/pos/modules/customerdisplay/models"
)

func TestDisplayConfig_Serialization(t *testing.T) {
	cfg := cd_models.DisplayConfig{
		Id:                "dsp-1",
		DisplayName:       "Lobby Screen",
		ShowPromotions:    true,
		ShowMenu:          true,
		ShowOrderStatus:   false,
		ShowWaitTime:      true,
		AutoSlideInterval: 10,
		Theme:             "dark",
		WelcomeMessage:    "Welcome!",
		Active:            true,
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded cd_models.DisplayConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.DisplayName != "Lobby Screen" {
		t.Errorf("expected DisplayName='Lobby Screen', got %s", decoded.DisplayName)
	}
	if !decoded.ShowPromotions {
		t.Error("expected ShowPromotions=true")
	}
	if decoded.AutoSlideInterval != 10 {
		t.Errorf("expected AutoSlideInterval=10, got %d", decoded.AutoSlideInterval)
	}
}

func TestDisplayConfig_DefaultValues(t *testing.T) {
	cfg := cd_models.DisplayConfig{}

	if cfg.Active {
		t.Error("expected Active=false by default")
	}
	if cfg.AutoSlideInterval != 0 {
		t.Errorf("expected AutoSlideInterval=0, got %d", cfg.AutoSlideInterval)
	}
}

func TestDisplayContent_Serialization(t *testing.T) {
	content := cd_models.DisplayContent{
		Items: []cd_models.DisplayItem{
			{Type: "promotions", Content: "data"},
			{Type: "menu", Content: []string{"cat1", "cat2"}},
		},
		Interval:   8,
		Theme:      "light",
		WelcomeMsg: "Hello",
	}

	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded cd_models.DisplayContent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(decoded.Items))
	}
	if decoded.Items[0].Type != "promotions" {
		t.Errorf("expected item 0 type='promotions', got %s", decoded.Items[0].Type)
	}
}

func TestDisplayContent_Empty(t *testing.T) {
	content := cd_models.DisplayContent{
		Items:    []cd_models.DisplayItem{},
		Interval: 0,
		Theme:    "",
	}

	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("failed to marshal empty: %v", err)
	}

	var decoded cd_models.DisplayContent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(decoded.Items))
	}
}
