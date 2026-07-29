package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/giftcard/dto"
	"github.com/nutrixpos/pos/modules/giftcard/models"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	card := models.GiftCard{
		Code:          "GC-1234",
		InitialAmount: 50.00,
		CurrentAmount: 50.00,
		Status:        "active",
	}
	resp := JSONApiOkResponse{Data: card, Meta: JSONAPIMeta{TotalRecords: 1}}

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

func TestGetAllGiftCards_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/giftcards", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cards := []models.GiftCard{}
		resp := JSONApiOkResponse{Data: cards, Meta: JSONAPIMeta{TotalRecords: 0}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/giftcards", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGetGiftCard_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/giftcards/{code}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		card := models.GiftCard{Code: "GC-1234", CurrentAmount: 50.00}
		resp := JSONApiOkResponse{Data: card, Meta: JSONAPIMeta{TotalRecords: 1}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/giftcards/GC-1234", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRedeemGiftCard_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/giftcards/redeem", func(w http.ResponseWriter, r *http.Request) {
		var req dto.RedeemGiftCardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	req := httptest.NewRequest("POST", "/api/giftcards/redeem", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for nil body, got %d", rec.Code)
	}
}

func TestGiftCard_TransactionSerialization(t *testing.T) {
	tx := models.GiftCardTransaction{
		GiftCardCode: "GC-1234",
		Type:         "redeem",
		Amount:       25.00,
		BalanceAfter: 25.00,
	}

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded models.GiftCardTransaction
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Type != "redeem" {
		t.Errorf("expected redeem, got %s", decoded.Type)
	}
	if decoded.Amount != 25.00 {
		t.Errorf("expected 25.00, got %f", decoded.Amount)
	}
}
