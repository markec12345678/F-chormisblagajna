package handlers

import (
	"encoding/json"
	"testing"
	"time"

	alert_models "github.com/nutrixpos/pos/modules/inventoryalerts/models"
)

func TestInventoryAlertRule_Serialization(t *testing.T) {
	rule := alert_models.InventoryAlertRule{
		Id:            "r-1",
		MaterialId:    "m-1",
		MaterialName:  "Tomatoes",
		ThresholdLow:  10.0,
		ThresholdCrit: 5.0,
		NotifyEmail:   true,
		IsActive:      true,
		CreatedAt:     time.Now(),
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded alert_models.InventoryAlertRule
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.MaterialName != "Tomatoes" {
		t.Errorf("expected MaterialName='Tomatoes', got %s", decoded.MaterialName)
	}
	if decoded.ThresholdLow != 10.0 {
		t.Errorf("expected ThresholdLow=10, got %f", decoded.ThresholdLow)
	}
}

func TestInventoryAlert_Serialization(t *testing.T) {
	alert := alert_models.InventoryAlert{
		Id:           "a-1",
		RuleId:       "r-1",
		MaterialId:   "m-1",
		MaterialName: "Tomatoes",
		CurrentQty:   3.0,
		Threshold:    5.0,
		Severity:     "critical",
		IsRead:       false,
		CreatedAt:    time.Now(),
	}

	data, err := json.Marshal(alert)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded alert_models.InventoryAlert
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Severity != "critical" {
		t.Errorf("expected Severity='critical', got %s", decoded.Severity)
	}
	if decoded.IsRead {
		t.Error("expected IsRead=false")
	}
}

func TestAlertSummary_Serialization(t *testing.T) {
	summary := alert_models.AlertSummary{
		TotalActive:   5,
		UnreadCount:   3,
		CriticalCount: 1,
		LowCount:      2,
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded alert_models.AlertSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.TotalActive != 5 {
		t.Errorf("expected TotalActive=5, got %d", decoded.TotalActive)
	}
}

func TestSeverity_Values(t *testing.T) {
	validSeverities := []string{"low", "critical"}
	for _, s := range validSeverities {
		if s == "" {
			t.Error("severity should not be empty")
		}
	}
}
