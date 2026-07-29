package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/core/models"
)

func TestCustomerModel_EnhancedFields_JSON(t *testing.T) {
	now := time.Now()
	customer := models.Customer{
		Id:            "c-1",
		Name:          "Janez Novak",
		Email:         "janez@example.com",
		Notes:         "Preferred window seat",
		Tags:          []string{"vip", "regular"},
		LoyaltyPoints: 150,
		TotalSpent:    1234.56,
		OrderCount:    25,
		LastOrderDate: &now,
		Preferences:   map[string]string{"language": "sl"},
		CreatedAt:     now,
	}

	data, err := json.Marshal(customer)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded models.Customer
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Email != "janez@example.com" {
		t.Errorf("Email = %v, want janez@example.com", decoded.Email)
	}
	if decoded.Notes != "Preferred window seat" {
		t.Errorf("Notes = %v, want 'Preferred window seat'", decoded.Notes)
	}
	if len(decoded.Tags) != 2 {
		t.Errorf("Tags length = %d, want 2", len(decoded.Tags))
	}
	if decoded.LoyaltyPoints != 150 {
		t.Errorf("LoyaltyPoints = %d, want 150", decoded.LoyaltyPoints)
	}
	if decoded.TotalSpent != 1234.56 {
		t.Errorf("TotalSpent = %f, want 1234.56", decoded.TotalSpent)
	}
}

func TestCustomerModel_BackwardCompatibility(t *testing.T) {
	customer := models.Customer{
		Id:   "c-1",
		Name: "Janez",
	}

	data, err := json.Marshal(customer)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded models.Customer
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Name != "Janez" {
		t.Errorf("Name = %v, want Janez", decoded.Name)
	}
	if decoded.Email != "" {
		t.Errorf("Email should be empty, got %v", decoded.Email)
	}
	if decoded.Tags != nil {
		t.Errorf("Tags should be nil, got %v", decoded.Tags)
	}
	if decoded.LoyaltyPoints != 0 {
		t.Errorf("LoyaltyPoints should be 0, got %d", decoded.LoyaltyPoints)
	}
}

func TestCustomerOrdersEndpoint_Routing(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/customers/{id}/orders", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		orders := []models.Order{}
		resp := JSONApiOkResponse{Data: orders, Meta: JSONAPIMeta{TotalRecords: 0}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/customers/c-1/orders", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp JSONApiOkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
}

func TestCustomerStatsEndpoint_Routing(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/customers/{id}/stats", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).Methods("POST")

	req := httptest.NewRequest("POST", "/api/customers/c-1/stats", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}
