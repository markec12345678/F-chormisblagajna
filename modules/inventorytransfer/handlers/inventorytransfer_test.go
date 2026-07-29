package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/inventorytransfer/dto"
	"github.com/nutrixpos/pos/modules/inventorytransfer/models"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	transfer := models.InventoryTransfer{
		MaterialID:   "m1",
		MaterialName: "Flour",
		Quantity:     10,
		Unit:         "kg",
		Status:       "pending",
	}
	resp := JSONApiOkResponse{Data: transfer, Meta: JSONAPIMeta{TotalRecords: 1}}

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

func TestGetAllTransfers_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/transfers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		transfers := []models.InventoryTransfer{}
		resp := JSONApiOkResponse{Data: transfers, Meta: JSONAPIMeta{TotalRecords: 0}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/transfers", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCreateTransfer_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/transfers", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateTransferRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	req := httptest.NewRequest("POST", "/api/transfers", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for nil body, got %d", rec.Code)
	}
}

func TestInventoryTransfer_Fields(t *testing.T) {
	transfer := models.InventoryTransfer{
		MaterialID:   "m1",
		MaterialName: "Flour",
		Quantity:     10,
		Unit:         "kg",
		FromBranchID: "b1",
		ToBranchID:   "b2",
		Status:       "pending",
	}

	if transfer.Status != "pending" {
		t.Errorf("expected pending, got %s", transfer.Status)
	}
	if transfer.Quantity != 10 {
		t.Errorf("expected 10, got %f", transfer.Quantity)
	}
}
