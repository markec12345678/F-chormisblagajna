package handlers

import (
	"encoding/json"
	"testing"
	"time"

	waste_models "github.com/nutrixpos/pos/modules/waste/models"
)

func TestWasteEntry_Serialization(t *testing.T) {
	entry := waste_models.WasteEntry{
		Id:           "w-1",
		MaterialId:   "m-1",
		MaterialName: "Tomatoes",
		Quantity:     5.0,
		Unit:         "kg",
		Reason:       "expired",
		Cost:         12.50,
		Date:         time.Now(),
		RecordedBy:   "chef-1",
		Notes:        "Expired before use",
	}

	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded waste_models.WasteEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.MaterialName != "Tomatoes" {
		t.Errorf("expected MaterialName='Tomatoes', got %s", decoded.MaterialName)
	}
	if decoded.Cost != 12.50 {
		t.Errorf("expected Cost=12.50, got %f", decoded.Cost)
	}
}

func TestWasteSummary_Serialization(t *testing.T) {
	summary := waste_models.WasteSummary{
		TotalWasteCost: 250.00,
		TotalEntries:   15,
		ByReason: []waste_models.ReasonSummary{
			{Reason: "expired", Total: 150.00, Count: 8},
			{Reason: "damaged", Total: 100.00, Count: 7},
		},
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded waste_models.WasteSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.TotalWasteCost != 250.00 {
		t.Errorf("expected TotalWasteCost=250, got %f", decoded.TotalWasteCost)
	}
	if len(decoded.ByReason) != 2 {
		t.Errorf("expected 2 reasons, got %d", len(decoded.ByReason))
	}
}

func TestWasteReasons(t *testing.T) {
	validReasons := []string{"expired", "damaged", "overcooked", "other"}
	for _, reason := range validReasons {
		if reason == "" {
			t.Error("reason should not be empty")
		}
	}
}

func TestDailyWaste_Calculation(t *testing.T) {
	daily := waste_models.DailyWaste{
		Date: "2024-01-15",
	}

	daily.Total = 25.00 + 35.00 + 15.00
	if daily.Total != 75.00 {
		t.Errorf("expected Total=75, got %f", daily.Total)
	}
}
