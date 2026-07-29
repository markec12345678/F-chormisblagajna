package handlers

import (
	"encoding/json"
	"testing"

	ly_models "github.com/nutrixpos/pos/modules/loyalty/models"
)

func TestLoyaltyCard_Serialization(t *testing.T) {
	card := ly_models.LoyaltyCard{
		Id:           "lc-1",
		CustomerId:   "cust-1",
		CustomerName: "John Doe",
		Points:       150,
		Tier:         "silver",
		TotalSpent:   500.00,
		VisitCount:   10,
		Active:       true,
	}

	data, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ly_models.LoyaltyCard
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Points != 150 {
		t.Errorf("expected Points=150, got %d", decoded.Points)
	}
	if decoded.Tier != "silver" {
		t.Errorf("expected Tier='silver', got %s", decoded.Tier)
	}
}

func TestReward_Serialization(t *testing.T) {
	reward := ly_models.Reward{
		Id:          "rw-1",
		Name:        "Free Coffee",
		Description: "Get a free coffee",
		PointsCost:  100,
		RewardType:  "free_item",
		FreeItemId:  "coffee-1",
		Active:      true,
	}

	data, err := json.Marshal(reward)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ly_models.Reward
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Name != "Free Coffee" {
		t.Errorf("expected Name='Free Coffee', got %s", decoded.Name)
	}
}

func TestRedemption_Serialization(t *testing.T) {
	red := ly_models.Redemption{
		Id:          "rd-1",
		CardId:      "lc-1",
		RewardId:    "rw-1",
		RewardName:  "Free Coffee",
		PointsSpent: 100,
		RedeemedAt:  "2024-01-15 10:30",
	}

	data, err := json.Marshal(red)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ly_models.Redemption
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.PointsSpent != 100 {
		t.Errorf("expected PointsSpent=100, got %d", decoded.PointsSpent)
	}
}

func TestLoyaltySettings_Serialization(t *testing.T) {
	settings := ly_models.LoyaltySettings{
		PointsPerEuro: 10,
		EuroPerPoint:  0.01,
		WelcomePoints: 50,
		TierThresholds: map[string]float64{"bronze": 0, "silver": 500},
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ly_models.LoyaltySettings
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.PointsPerEuro != 10 {
		t.Errorf("expected PointsPerEuro=10, got %f", decoded.PointsPerEuro)
	}
}
