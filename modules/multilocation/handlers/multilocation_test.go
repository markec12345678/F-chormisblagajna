package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/multilocation/models"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	dashboard := models.LocationDashboard{
		TotalRevenue:  1500.50,
		TotalOrders:   45,
		TotalBranches: 3,
		AvgOrderValue: 33.34,
	}
	resp := JSONApiOkResponse{Data: dashboard, Meta: JSONAPIMeta{TotalRecords: 1}}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded JSONApiOkResponse
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Meta.TotalRecords != 1 {
		t.Errorf("Expected TotalRecords=1, got %d", decoded.Meta.TotalRecords)
	}
}

func TestGetDashboard_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/multilocation/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dashboard := models.LocationDashboard{TotalBranches: 0, Branches: make([]models.BranchStats, 0)}
		resp := JSONApiOkResponse{Data: dashboard, Meta: JSONAPIMeta{TotalRecords: 1}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/multilocation/dashboard", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGetComparison_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/multilocation/comparison", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		comparison := []models.BranchComparison{}
		resp := JSONApiOkResponse{Data: comparison, Meta: JSONAPIMeta{TotalRecords: 0}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/multilocation/comparison", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestJSONAPIMeta_Serialization(t *testing.T) {
	meta := JSONAPIMeta{TotalRecords: 5}
	jsonBytes, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("Failed to marshal meta: %v", err)
	}

	var decoded JSONAPIMeta
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal meta: %v", err)
	}

	if decoded.TotalRecords != 5 {
		t.Errorf("Expected TotalRecords=5, got %d", decoded.TotalRecords)
	}
}

func TestLocationDashboard_Fields(t *testing.T) {
	dashboard := models.LocationDashboard{
		TotalRevenue:  5000,
		TotalOrders:   100,
		TotalBranches: 2,
		AvgOrderValue: 50,
		Branches: []models.BranchStats{
			{BranchID: "b1", BranchName: "Main", TodayRevenue: 3000, TodayOrders: 60},
			{BranchID: "b2", BranchName: "North", TodayRevenue: 2000, TodayOrders: 40},
		},
	}

	if dashboard.TotalBranches != 2 {
		t.Errorf("expected 2 branches, got %d", dashboard.TotalBranches)
	}
	if len(dashboard.Branches) != 2 {
		t.Errorf("expected 2 branch stats, got %d", len(dashboard.Branches))
	}
	if dashboard.Branches[0].BranchName != "Main" {
		t.Errorf("expected Main, got %s", dashboard.Branches[0].BranchName)
	}
}
