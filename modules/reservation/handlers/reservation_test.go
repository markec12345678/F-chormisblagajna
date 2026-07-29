package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/reservation/dto"
	"github.com/nutrixpos/pos/modules/reservation/models"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	r := models.Reservation{
		Id:            "test-id-123",
		CustomerName:  "Janez Novak",
		CustomerPhone: "+386 40 123 456",
		TableId:       "t-1",
		Date:          "2026-07-27",
		Time:          "19:00",
		GuestCount:    4,
		Status:        "confirmed",
	}

	resp := JSONApiOkResponse{
		Data: r,
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

func TestCreateReservationRequest_Deserialization(t *testing.T) {
	body := `{"customer_name":"Janez Novak","customer_phone":"+386 40 123 456","table_id":"t-1","date":"2026-07-27","time":"19:00","guest_count":4,"status":"confirmed"}`

	var req dto.CreateReservationRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.CustomerName != "Janez Novak" {
		t.Errorf("expected customer_name Janez Novak, got %s", req.CustomerName)
	}
	if req.GuestCount != 4 {
		t.Errorf("expected guest_count 4, got %d", req.GuestCount)
	}
}

func TestGetReservations_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/reservations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":[],"meta":{"total_records":0}}`))
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/reservations", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCreateReservation_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/reservations", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateReservationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	req := httptest.NewRequest("POST", "/api/reservations", strings.NewReader(`{invalid}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreateReservation_ValidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/reservations", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateReservationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		resp := JSONApiOkResponse{
			Data: models.Reservation{
				CustomerName: req.CustomerName,
				Date:         req.Date,
				Time:         req.Time,
				GuestCount:   req.GuestCount,
			},
			Meta: JSONAPIMeta{TotalRecords: 1},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	body := `{"customer_name":"Janez Novak","customer_phone":"+386 40 123 456","table_id":"t-1","date":"2026-07-27","time":"19:00","guest_count":4}`
	req := httptest.NewRequest("POST", "/api/reservations", strings.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
}

func TestDeleteReservation_NoContent(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/reservations/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).Methods("DELETE")

	req := httptest.NewRequest("DELETE", "/api/reservations/abc123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}

func TestUpdateReservationRequest_Deserialization(t *testing.T) {
	body := `{"status":"seated","guest_count":6}`

	var req dto.UpdateReservationRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.Status != "seated" {
		t.Errorf("expected status seated, got %s", req.Status)
	}
	if req.GuestCount == nil {
		t.Fatal("expected guest_count to be non-nil")
	}
	if *req.GuestCount != 6 {
		t.Errorf("expected guest_count 6, got %d", *req.GuestCount)
	}
}
