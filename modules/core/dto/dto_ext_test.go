package dto

import (
	"encoding/json"
	"testing"
)

func TestRecipeAvailability_JSON(t *testing.T) {
	r := RecipeAvailability{
		RecipeId:  "recipe-1",
		Available: 10.5,
		Ready:     5.0,
		ComponentRequirements: map[string]float64{
			"mat-1": 2.0,
			"mat-2": 0.5,
		},
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded RecipeAvailability
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.RecipeId != r.RecipeId {
		t.Errorf("RecipeId = %v, want %v", decoded.RecipeId, r.RecipeId)
	}
	if decoded.Available != r.Available {
		t.Errorf("Available = %v, want %v", decoded.Available, r.Available)
	}
	if decoded.Ready != r.Ready {
		t.Errorf("Ready = %v, want %v", decoded.Ready, r.Ready)
	}
	if len(decoded.ComponentRequirements) != 2 {
		t.Errorf("ComponentRequirements length = %v, want 2", len(decoded.ComponentRequirements))
	}
	if decoded.ComponentRequirements["mat-1"] != 2.0 {
		t.Errorf("ComponentRequirements[mat-1] = %v, want 2.0", decoded.ComponentRequirements["mat-1"])
	}
}

func TestRecipeAvailability_Empty(t *testing.T) {
	r := RecipeAvailability{
		RecipeId:              "recipe-2",
		Available:             0,
		Ready:                 0,
		ComponentRequirements: map[string]float64{},
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded RecipeAvailability
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.ComponentRequirements == nil {
		t.Error("ComponentRequirements should not be nil after unmarshal")
	}
}

func TestComponentQuantity_JSON(t *testing.T) {
	c := ComponentQuantity{
		ComponentId: "comp-1",
		Quantity:    5.5,
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded ComponentQuantity
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.ComponentId != c.ComponentId {
		t.Errorf("ComponentId = %v, want %v", decoded.ComponentId, c.ComponentId)
	}
	if decoded.Quantity != c.Quantity {
		t.Errorf("Quantity = %v, want %v", decoded.Quantity, c.Quantity)
	}
}

func TestGetComponentConsumeLogsRequest_JSON(t *testing.T) {
	r := GetComponentConsumeLogsRequest{Name: "Flour"}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded GetComponentConsumeLogsRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Name != r.Name {
		t.Errorf("Name = %v, want %v", decoded.Name, r.Name)
	}
}
