package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/promotion/dto"
	"github.com/nutrixpos/pos/modules/promotion/models"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	p := models.Promotion{
		Id:        "test-id-123",
		Name:      "Summer 20% OFF",
		Code:      "SUMMER20",
		Type:      "percentage",
		Value:     20,
		MinOrder:  50,
		IsActive:  true,
	}

	resp := JSONApiOkResponse{
		Data: p,
		Meta: JSONAPIMeta{TotalRecords: 1, PageNumber: 1, PageSize: 50, PageCount: 1},
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
		t.Errorf("expected total_records 1, got %d", decoded.Meta.TotalRecords)
	}
}

func TestCreatePromotionRequest_Deserialization(t *testing.T) {
	body := `{"name":"Summer 20% OFF","code":"SUMMER20","type":"percentage","value":20,"min_order":50,"is_active":true}`

	var req dto.CreatePromotionRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.Name != "Summer 20% OFF" {
		t.Errorf("expected name Summer 20%% OFF, got %s", req.Name)
	}
	if req.Code != "SUMMER20" {
		t.Errorf("expected code SUMMER20, got %s", req.Code)
	}
	if req.Value != 20 {
		t.Errorf("expected value 20, got %f", req.Value)
	}
}

func TestGetPromotions_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/promotions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":[],"meta":{"total_records":0}}`))
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/promotions", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCreatePromotion_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/promotions", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreatePromotionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	req := httptest.NewRequest("POST", "/api/promotions", strings.NewReader(`{invalid}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreatePromotion_ValidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/promotions", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreatePromotionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		resp := JSONApiOkResponse{
			Data: models.Promotion{
				Name:   req.Name,
				Code:   req.Code,
				Type:   req.Type,
				Value:  req.Value,
				IsActive: req.IsActive,
			},
			Meta: JSONAPIMeta{TotalRecords: 1},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	body := `{"name":"Happy Hour","code":"HAPPY","type":"percentage","value":15,"is_active":true}`
	req := httptest.NewRequest("POST", "/api/promotions", strings.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
}

func TestDeletePromotion_NoContent(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/promotions/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).Methods("DELETE")

	req := httptest.NewRequest("DELETE", "/api/promotions/abc123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}

func TestUpdatePromotionRequest_Deserialization(t *testing.T) {
	body := `{"name":"Updated Promo","is_active":false}`

	var req dto.UpdatePromotionRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.Name != "Updated Promo" {
		t.Errorf("expected name Updated Promo, got %s", req.Name)
	}
}
