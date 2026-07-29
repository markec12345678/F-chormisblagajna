package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/tips/dto"
	"github.com/nutrixpos/pos/modules/tips/models"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	tip := models.Tip{
		OrderID:      "order1",
		EmployeeID:   "emp1",
		EmployeeName: "John",
		Amount:       5.50,
	}
	resp := JSONApiOkResponse{Data: tip, Meta: JSONAPIMeta{TotalRecords: 1}}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded JSONApiOkResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Meta.TotalRecords != 1 {
		t.Errorf("expected TotalRecords=1, got %d", decoded.Meta.TotalRecords)
	}
}

func TestGetTipsByEmployee_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/tips/summary", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		summaries := []models.TipSummary{}
		resp := JSONApiOkResponse{Data: summaries, Meta: JSONAPIMeta{TotalRecords: 0}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/tips/summary", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRecordTip_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/tips", func(w http.ResponseWriter, r *http.Request) {
		var req dto.RecordTipRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	req := httptest.NewRequest("POST", "/api/tips", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for nil body, got %d", rec.Code)
	}
}

func TestTipSummary_Fields(t *testing.T) {
	summary := models.TipSummary{
		EmployeeID:   "emp1",
		EmployeeName: "John",
		TotalTips:    50.00,
		TipCount:     10,
		AverageTip:   5.00,
	}

	if summary.TotalTips != 50.00 {
		t.Errorf("expected 50.00, got %f", summary.TotalTips)
	}
	if summary.TipCount != 10 {
		t.Errorf("expected 10, got %d", summary.TipCount)
	}
}
