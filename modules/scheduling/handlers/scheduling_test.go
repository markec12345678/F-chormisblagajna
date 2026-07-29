package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/scheduling/dto"
	"github.com/nutrixpos/pos/modules/scheduling/models"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	shift := models.Shift{
		Id:         "test-id-123",
		EmployeeId: "emp-1",
		BranchId:   "br-1",
		Date:       "2026-07-27",
		StartTime:  "08:00",
		EndTime:    "16:00",
		Role:       "cashier",
		Status:     "scheduled",
	}

	resp := JSONApiOkResponse{
		Data: shift,
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

func TestCreateShiftRequest_Deserialization(t *testing.T) {
	body := `{"employee_id":"emp-1","branch_id":"br-1","date":"2026-07-27","start_time":"08:00","end_time":"16:00","role":"cashier","status":"scheduled","notes":"Morning shift"}`

	var req dto.CreateShiftRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.EmployeeId != "emp-1" {
		t.Errorf("expected employee_id emp-1, got %s", req.EmployeeId)
	}
	if req.Date != "2026-07-27" {
		t.Errorf("expected date 2026-07-27, got %s", req.Date)
	}
	if req.Role != "cashier" {
		t.Errorf("expected role cashier, got %s", req.Role)
	}
}

func TestGetShifts_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/shifts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":[],"meta":{"total_records":0}}`))
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/shifts", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCreateShift_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/shifts", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateShiftRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	req := httptest.NewRequest("POST", "/api/shifts", strings.NewReader(`{invalid}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreateShift_ValidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/shifts", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateShiftRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		resp := JSONApiOkResponse{
			Data: models.Shift{
				EmployeeId: req.EmployeeId,
				Date:       req.Date,
				Role:       req.Role,
				Status:     req.Status,
			},
			Meta: JSONAPIMeta{TotalRecords: 1},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	body := `{"employee_id":"emp-1","branch_id":"br-1","date":"2026-07-27","start_time":"08:00","end_time":"16:00","role":"cashier","status":"scheduled"}`
	req := httptest.NewRequest("POST", "/api/shifts", strings.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}

	var resp JSONApiOkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
}

func TestDeleteShift_NoContent(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/shifts/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).Methods("DELETE")

	req := httptest.NewRequest("DELETE", "/api/shifts/abc123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}

func TestUpdateShiftRequest_Deserialization(t *testing.T) {
	body := `{"status":"confirmed","notes":"Updated notes"}`

	var req dto.UpdateShiftRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.Status != "confirmed" {
		t.Errorf("expected status confirmed, got %s", req.Status)
	}
	if req.Notes != "Updated notes" {
		t.Errorf("expected notes 'Updated notes', got %s", req.Notes)
	}
}
