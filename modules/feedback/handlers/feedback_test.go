package handlers

import (
	"encoding/json"
	"testing"
	"time"

	fb_models "github.com/nutrixpos/pos/modules/feedback/models"
)

func TestFeedback_Serialization(t *testing.T) {
	fb := fb_models.Feedback{
		Id:         "fb-1",
		OrderId:    "o-1",
		Rating:     5,
		Comment:    "Excellent food!",
		Category:   "food",
		Anonymous:  false,
		CreatedAt:  time.Now(),
		Responded:  false,
	}

	data, err := json.Marshal(fb)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded fb_models.Feedback
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Rating != 5 {
		t.Errorf("expected Rating=5, got %d", decoded.Rating)
	}
	if decoded.Category != "food" {
		t.Errorf("expected Category='food', got %s", decoded.Category)
	}
}

func TestFeedback_WithResponse(t *testing.T) {
	fb := fb_models.Feedback{
		Id:        "fb-2",
		Rating:    3,
		Comment:   "OK service",
		Category:  "service",
		Responded: true,
		Response:  "Thank you for your feedback!",
	}

	data, err := json.Marshal(fb)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded fb_models.Feedback
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !decoded.Responded {
		t.Error("expected Responded=true")
	}
	if decoded.Response != "Thank you for your feedback!" {
		t.Errorf("unexpected response: %s", decoded.Response)
	}
}

func TestFeedback_Anonymous(t *testing.T) {
	fb := fb_models.Feedback{
		Id:        "fb-3",
		Rating:    4,
		Anonymous: true,
	}

	data, err := json.Marshal(fb)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded fb_models.Feedback
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !decoded.Anonymous {
		t.Error("expected Anonymous=true")
	}
}

func TestFeedbackSummary_Serialization(t *testing.T) {
	summary := fb_models.FeedbackSummary{
		TotalFeedbacks: 10,
		AverageRating:  4.2,
		RatingDist:     map[int]int{5: 5, 4: 3, 3: 2},
		CategoryAvg:    map[string]float64{"food": 4.5, "service": 3.8},
		RecentFeedbacks: []fb_models.Feedback{
			{Id: "fb-1", Rating: 5},
			{Id: "fb-2", Rating: 4},
		},
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded fb_models.FeedbackSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.TotalFeedbacks != 10 {
		t.Errorf("expected TotalFeedbacks=10, got %d", decoded.TotalFeedbacks)
	}
}
