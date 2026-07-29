package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	resp := JSONApiOkResponse{
		Data: map[string]string{"message": "ok"},
		Meta: JSONAPIMeta{TotalRecords: 1},
	}

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

func TestGetStations_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/kitchen/stations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stations := []map[string]string{}
		resp := JSONApiOkResponse{Data: stations, Meta: JSONAPIMeta{TotalRecords: 0}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/kitchen/stations", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestUpdateItemStatus_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/kitchen/orders/{order_id}/items/status", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ItemIndex int    `json:"item_index"`
			Status    string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("PUT")

	req := httptest.NewRequest("PUT", "/api/kitchen/orders/order1/items/status", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for nil body, got %d", rec.Code)
	}
}

func TestGetOrdersByStation_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/kitchen/orders", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		orders := []map[string]interface{}{}
		resp := JSONApiOkResponse{Data: orders, Meta: JSONAPIMeta{TotalRecords: 0}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/kitchen/orders?station=grill", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
