package handlers

import (
	"encoding/json"
	"testing"
	"time"

	tr_models "github.com/nutrixpos/pos/modules/training/models"
)

func TestTrainingSession_Serialization(t *testing.T) {
	session := tr_models.TrainingSession{
		Id:         "tr-1",
		UserId:     "user-1",
		Module:     "cashier",
		StartedAt:  time.Now(),
		Score:      40,
		MaxScore:   100,
		Completed:  false,
		StepsDone:  2,
		TotalSteps: 5,
	}

	data, err := json.Marshal(session)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded tr_models.TrainingSession
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Module != "cashier" {
		t.Errorf("expected Module='cashier', got %s", decoded.Module)
	}
	if decoded.StepsDone != 2 {
		t.Errorf("expected StepsDone=2, got %d", decoded.StepsDone)
	}
}

func TestTrainingModule_Serialization(t *testing.T) {
	mod := tr_models.TrainingModule{
		Key:         "kitchen",
		Name:        "Kitchen Training",
		Description: "Learn kitchen workflow",
		Icon:        "pi pi-fire",
		Steps:       4,
	}

	data, err := json.Marshal(mod)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded tr_models.TrainingModule
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Name != "Kitchen Training" {
		t.Errorf("expected Name='Kitchen Training', got %s", decoded.Name)
	}
}

func TestTrainingStep_Serialization(t *testing.T) {
	step := tr_models.TrainingStep{
		Id:          "cs-1",
		Module:      "cashier",
		Title:       "Open Order",
		Description: "Click new order",
		Action:      "navigate_new_order",
		ExpectedOut: "Screen opens",
	}

	data, err := json.Marshal(step)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded tr_models.TrainingStep
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Title != "Open Order" {
		t.Errorf("expected Title='Open Order', got %s", decoded.Title)
	}
}

func TestTrainingProgress_Serialization(t *testing.T) {
	prog := tr_models.TrainingProgress{
		SessionId:     "tr-1",
		Module:        "cashier",
		StartedAt:     "2024-01-15 10:00",
		StepsDone:     3,
		TotalSteps:    5,
		Score:         60,
		MaxScore:      100,
		CompletionPct: 60.0,
		Completed:     false,
	}

	data, err := json.Marshal(prog)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded tr_models.TrainingProgress
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.CompletionPct != 60.0 {
		t.Errorf("expected CompletionPct=60, got %f", decoded.CompletionPct)
	}
}
